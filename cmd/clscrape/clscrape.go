package main

import (
    "fmt"
    //"encoding/xml"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "os"

    "github.com/gorilla/Schema"
    "github.com/jessevdk/go-flags"
    //"github.com/PuerkitoBio/goquery"
    "gopkg.in/yaml.v2"
)

type Config struct {
    QueryURL                string
    Search struct {
        Scheme              string  `yaml:"Scheme"`
        Location            string  `yaml:"Location"`
        URL                 string  `yaml:"URL"`
    } `yaml:"Search"`

    SMTP struct {
        Host                string  `yaml:"host,omitempty"`
        Port                string  `yaml:"port,omitempty"`
        User                string  `yaml:"user,omitempty"`
        Pass                string  `yaml:"pass,omitempty"`
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

type Listing struct {
    Title       string `xml:"title"`
    Description string
    Link        string
    Date        string
    Image       string
}

var opts struct {
    File string `short:"i" long:"input" description:"Yaml-formatted configuration file" required:"true"`
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
    fmt.Println(c.QueryURL)

    res, err := http.Get(c.QueryURL)
    defer res.Body.Close()

    html, err := ioutil.ReadAll(res.Body)

    if err != nil {
        os.Exit(1)
    }
    if res.StatusCode != 200 {
        log.Fatalf("Status error code: %d %s", res.StatusCode, res.Status)
    }

    xmlDecode(string(html))
}

func xmlDecode(res string) {
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