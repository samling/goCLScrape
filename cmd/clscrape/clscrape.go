package main

import (
    //"bytes"
    //"crypto/tls"
    "encoding/json"
    "fmt"
    //"html/template"
    "io/ioutil"
    "log"
    "net/http"
    "net/smtp"
    "net/url"
    "os"
    "regexp"
    "strings"

    "github.com/gorilla/Schema"
    "github.com/jessevdk/go-flags"
    "github.com/PuerkitoBio/goquery"
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
        Host                string   `yaml:"host,omitempty"`
        Port                string   `yaml:"port,omitempty"`
        User                string   `yaml:"user,omitempty"`
        Pass                string   `yaml:"pass,omitempty"`
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
    Price       string      `json:"price"`
    Location    string      `json:"location"`
    Link        string      `json:"link"`
}

type Request struct {
    from        string
    to          []string
    subject     string
    body        string
}

var opts struct {
    File        string `short:"i" long:"input" description:"Yaml-formatted configuration file" required:"true"`
}

var auth smtp.Auth

func main() {
    args := os.Args
    args, err := flags.ParseArgs(&opts, args)
    if err != nil {
        return
    }
    configFile := opts.File

    c := Config{}
    c.getConf(configFile)

    // Connect to gmail SMTP
    auth = smtp.PlainAuth("", c.SMTP.Host, c.SMTP.User, c.SMTP.Pass)

    // Send a request
    r := NewRequest([]string{"samlingx@gmail.com"}, "Hello world", "Hello world")
    fmt.Printf("%v", r)
    ok, _ := r.SendEmail(c)
    fmt.Println(ok)

    listings := Listings{}
    listings.getAll(c.QueryURL, c.Search.Filter)

    listings_json, err := json.MarshalIndent(listings.Listings, "", "    ")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s\n\n", string(listings_json))
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

	doc.Find("p.result-info").Each(func(i int, s *goquery.Selection) {
        match := Listing{}

        title := s.Find(".result-title").Text()
        price := s.Find(".result-meta > .result-price").Text()
        location := s.Find(".result-meta > .result-hood").Text()
        link, _ := s.Find("a").Attr("href")

        filterRegex := strings.Join(filterList, "|")
        r := regexp.MustCompile(filterRegex)
        var matches []string
        var replacer = strings.NewReplacer("(", "", ")", "")
        titleMatches := r.FindAllString(strings.ToUpper(title), -1)
        locMatches := r.FindAllString(strings.ToUpper(location), -1)
        matches = append(matches, locMatches...)
        matches = append(matches, titleMatches...)

        if matches == nil {
           match.Title = title
           match.Price = price
           match.Location = replacer.Replace(strings.TrimSpace(location))
           match.Link = link
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

func NewRequest(to []string, subject string, body string) *Request {
    return &Request{
        to:         to,
        subject:    subject,
        body:       body,
    }
}

func (r *Request) SendEmail(c Config) (bool, error) {
    mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
    subject := "Subject: " + r.subject + "!\n"
    msg := []byte(subject + mime + "\n" + r.body)
    addr := c.SMTP.Host + ":" + c.SMTP.Port

    if err := smtp.SendMail(addr, auth, "samlingx@gmail.com", r.to, msg); err != nil {
        fmt.Printf("SMTP Error: %s", err)
        return false, err
    }
    return true, nil
}
