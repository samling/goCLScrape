package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/url"

    "github.com/gorilla/Schema"
    "gopkg.in/yaml.v2"
)

type Config struct {
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

func main() {
    c := Config{}
    c.getConf()

    q := createQueryString(c)
    fmt.Println(q)
}

func (c *Config) getConf() *Config {
    yamlFile, err := ioutil.ReadFile("Config.yaml")
    if err != nil {
        log.Printf("yamlFile.Get err #%s ", err)
    }

    err = yaml.Unmarshal(yamlFile, c)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return c
}

func createQueryString(c Config) *url.URL {
    u := new(url.URL)
    host := c.Search.Location + "." + c.Search.URL
    fmt.Println(host)
    form := url.Values{}

    encoder := schema.NewEncoder()
    encoder.SetAliasTag("yaml")
    encoder.Encode(c.Query, form)

    u.RawQuery = form.Encode()
    u.Scheme = c.Search.Scheme
    u.Path = host


    return u
}
