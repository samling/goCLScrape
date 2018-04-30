package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/url"

    "gopkg.in/yaml.v2"
)

func main() {
    m := getMap()
    fmt.Println(m["Search"]["URL"])
    q := createQueryString(m)
    fmt.Println(q)
}

func getMap() map[string]map[string]string {
    // Since the nested values are also a map, need to do some mapception
    m := make(map[string]map[string]string)

    yamlFile, err := ioutil.ReadFile("Config.yaml")
    if err != nil {
        log.Printf("yamlFile.Get err #%s ", err)
    }

    // Dynamically unmarshal key/value pairs into our map of string maps
    err = yaml.Unmarshal(yamlFile, m)
    if err != nil {
        log.Fatal(err)
    }

    return m
}

func createQueryString(m map[string]map[string]string) string {
    u, err := url.Parse(m["Search"]["URL"])
    if err != nil {
        log.Fatal(err)
    }

    q := u.Query()
    for k, v := range m["Query"] {
        q.Set(k, v)
        fmt.Printf("key[%v] value[%v]\n", k, v)
    }
    fmt.Println(q)
    u.RawQuery = q.Encode()
    fmt.Println(u)
    return m["Query"]["hasPic"]
}
