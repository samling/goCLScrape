package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/url"
    //"strings"

    "gopkg.in/yaml.v2"
)

type Config struct {
    Search struct {
        URL     string  `yaml:"URL"`
    } `yaml:"Search"`

    SMTP struct {
        Host    string  `yaml:"host,omitempty"`
        Port    string  `yaml:"port,omitempty"`
        User    string  `yaml:"user,omitempty"`
        Pass    string  `yaml:"pass,omitempty"`
    } `yaml:"SMTP,omitempty"`

    Query struct {
        HasPic  string  `yaml:"hasPic"`
    } `yaml:"Query,omitempty"`
}

func main() {
    c := Config{}
    c.getConf()
    fmt.Println(c)
    //m := getConfigMap()
    createQueryString(c)
    //fmt.Println(q)
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
    fmt.Printf("%v", c)

    return c
}

//func getConfigMap() map[string]map[string]string {
//    // Since the nested values are also a map, need to do some mapception
//    m := make(map[string]map[string]string)
//
//    yamlFile, err := ioutil.ReadFile("Config.yaml")
//    if err != nil {
//        log.Printf("yamlFile.Get err #%s ", err)
//    }
//
//    // Dynamically unmarshal key/value pairs into our map of string maps
//    err = yaml.Unmarshal(yamlFile, m)
//    if err != nil {
//        log.Fatal(err)
//    }
//
//    return m
//}

func createQueryString(c Config) *url.URL {
    u, err := url.Parse(c.Search.URL)
    if err != nil {
        log.Fatal(err)
    }

    // Create an empty query string that we Set() key/value pairs on; results in "key=value&key=value&..."
//    q := u.Query()
//    for k, v := range m["Query"] {
//        // TODO: Handle sequence of values so that they result in a series of duplicate keys, e.g. housing=1&housing=2&housing=3
//        if strings.Contains(v, ",") { // Tried giving values in YAML file like: "housing: 1,2"
//            s := strings.Split(v, ",")
//            for i, j := range s {
//                q.Set(k, j)
//                fmt.Println(i)
//                fmt.Println(j)
//                // Guess there'd be a 'continue' here or something
//            }
//            fmt.Println(k)
//            fmt.Println(v)
//        }
//        q.Set(k, v) // But what to do about this if it's a sequence...
//        // Would rather set value in config like:
//        // housing: [1 2]
//        // Problem is map explicitly contains strings, so results in: "line 63: cannot unmarshal !!seq into string"
//        // Tried changing map to map[string]map[string]interface{} but then it gets messy with type conversions etc.
//    }
//
//    u.RawQuery = q.Encode()
//
    return u
}
