package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "os"
    "regexp"
    "strings"
    "time"

    "github.com/gorilla/Schema"
    "github.com/jessevdk/go-flags"
    "github.com/PuerkitoBio/goquery"
    "gopkg.in/gomail.v2"
    "gopkg.in/yaml.v2"
)

type Config struct {
    QueryURL                string
    Search struct {
        Scheme              string   `yaml:"Scheme"`
        Location            string   `yaml:"Location"`
        URL                 string   `yaml:"URL"`
        Filter              []string `yaml:"Filter"`
    } `yaml:"Search"`

    SMTP struct {
        Name                string   `yaml:"name,omitempty"`
        Host                string   `yaml:"host,omitempty"`
        Port                int      `yaml:"port,omitempty"`
        User                string   `yaml:"user,omitempty"`
        Pass                string   `yaml:"pass,omitempty"`
        From                string   `yaml:"from,omitempty"`
        To                  string   `yaml:"to,omitempty"`
    } `yaml:"SMTP"`

    Query struct {
        Format              string   `yaml:"format,omitempty"`
        HasPic              string   `yaml:"hasPic,omitempty"`
        SrchType            string   `yaml:"srchType,omitempty"`
        BundleDuplicates    string   `yaml:"bundleDuplicates,omitempty"`
        MinPrice            string   `yaml:"min_price,omitempty"`
        MaxPrice            string   `yaml:"max_price,omitempty"`
        PostedToday         string   `yaml:"postedToday,omitempty"`
        SaleDate            string   `yaml:"sale_date,omitempty"`
        AvailabilityMode    string   `yaml:"availabilityMode,omitempty"`
        SearchDistance      string   `yaml:"search_distance,omitempty"`
        Postal              string   `yaml:"postal,omitempty"`
        SearchNearby        string   `yaml:"searchNearby,omitempty"`
        NearbyAreas         []string `yaml:"nearbyAreas,omitempty"`
        MinBedrooms         string   `yaml:"min_bedrooms,omitempty"`
        MaxBedrooms         string   `yaml:"max_bedrooms,omitempty"`
        MinBathrooms        string   `yaml:"min_bathrooms,omitempty"`
        MaxBathrooms        string   `yaml:"max_bathrooms,omitempty"`
        MinSqft             string   `yaml:"minSqft,omitempty"`
        MaxSqft             string   `yaml:"maxSqft,omitempty"`
        PetsCat             string   `yaml:"pets_cat,omitempty"`
        PetsDog             string   `yaml:"pets_dog,omitempty"`
        IsFurnished         string   `yaml:"is_furnished,omitempty"`
        Wheelchaccess       string   `yaml:"wheelchaccess,omitempty"`
        HousingType         []string `yaml:"housing_type,omitempty"`
        Laundry             []string `yaml:"laundry,omitempty"`
        Parking             []string `yaml:"parking,omitempty"`
    } `yaml:"Query"`
}

type Listings struct {
    Listings    []Listing   `json:"listings"`
}

type Listing struct {
    Title       string      `json:"title"`
    Date        string      `json:"date"`
    Price       string      `json:"price"`
    Location    string      `json:"location"`
    Link        string      `json:"link"`
    Image       string      `json:"image"`
}

var opts struct {
    File        string `short:"i" long:"input" description:"Yaml-formatted configuration file" required:"true"`
}

func main() {
    args := os.Args
    args, err := flags.ParseArgs(&opts, args)
    if err != nil {
        return
    }
    configFile := opts.File

    c := Config{}
    c.getConf(configFile)

    listings := Listings{}
    listings.getAll(c.QueryURL, c.Search.Filter)

    listings_json, err := json.MarshalIndent(listings.Listings, "", "    ")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s\n\n", string(listings_json))

    sendResults(c, listings)
}

func (listings *Listings) getAll(url string, filterList []string) {
    res, err := http.Get(url)
    if err != nil {
        log.Fatal("Unable to fetch URL")
    }
    defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("li.result-row").Each(func(i int, s *goquery.Selection) {
        match := Listing{}

        var images string
        var imageids []string
        title := s.Find("p.result-info > .result-title").Text()
        date := s.Find("p.result-info > .result-date").Text()
        price := s.Find("p.result-info > .result-meta > .result-price").Text()
        location := s.Find("p.result-info > .result-meta > .result-hood").Text()
        link, _ := s.Find("p.result-info > a").Attr("href")
        images, _ = s.Find("a.result-image").Attr("data-ids")

        // Filter matches using filter list in config
        filterRegex := strings.Join(filterList, "|")
        r := regexp.MustCompile(filterRegex)
        var matches []string
        var replacer = strings.NewReplacer("(", "", ")", "")
        titleMatches := r.FindAllString(strings.ToUpper(title), -1)
        locMatches := r.FindAllString(strings.ToUpper(location), -1)
        matches = append(matches, locMatches...)
        matches = append(matches, titleMatches...)

        // Select image from data-ids for email preview
        imageids = strings.Split(images, ",")
        image := strings.Replace(imageids[0], "1:", "", -1)

        // Add valid match to struct, append struct to slice of structs
        if matches == nil {
           match.Title = title
           match.Date = date
           match.Price = price
           match.Location = replacer.Replace(strings.TrimSpace(location))
           match.Link = link
           match.Image = image
           listings.Listings = append(listings.Listings, match)
        }
	})
}

func (c *Config) getConf(configFile string) *Config {
    yamlFile, err := ioutil.ReadFile(configFile)
    if err != nil {
        log.Printf("yamlFile.Get err #%s ", err)
    }

    err = yaml.Unmarshal(yamlFile, c)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }
    fmt.Printf("Filters: %v\n\n", c.Search.Filter)

    c.QueryURL = c.getURL()

    return c
}

func (c *Config) getURL() string {
    u := new(url.URL)

    host := c.Search.Location + "." + c.Search.URL

    u.Scheme = c.Search.Scheme
    u.Path = host

    form := url.Values{}

    encoder := schema.NewEncoder()
    encoder.SetAliasTag("yaml")
    encoder.Encode(c.Query, form)

    u.RawQuery = form.Encode()

    return u.String()
}

func sendResults(c Config, listings Listings) error {
    var body string
    body = body + "Hello " + c.SMTP.Name + ", here are your latest search results:<br>"
    for _, listing := range listings.Listings {
        body = body + "<h3>" + "<a href='" + listing.Link + "'>" + listing.Title + "</a>" + " - " + listing.Price + " "
        if len(listing.Location) > 0 {
            body = body + " | " + listing.Location + " "
        }
        body = body + " | <i>" + listing.Date + "</i></h3>"
        body = body + "<img src='https://images.craigslist.org/" + listing.Image + "_300x300.jpg'><br>"
    }
    body = body + "<br><br><br><a href='" + c.QueryURL + "'>Browse the results</a>"

    mail := gomail.NewMessage()
    mail.SetHeader("From", c.SMTP.From)
    mail.SetHeader("To", c.SMTP.To)
    mail.SetHeader("Subject", "Craigslist search results - " + time.Now().Format(time.RFC850))
    mail.SetBody("text/html", body)

    d := gomail.NewDialer(c.SMTP.Host, c.SMTP.Port, c.SMTP.User, c.SMTP.Pass)

    if err := d.DialAndSend(mail); err != nil {
        return err
    }
    return nil
}
