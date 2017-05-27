package main

import (
    "fmt"
    "encoding/json"
    "log"
)

type Lang struct {
    Name string
    Year int
    URL string
}

func main() {
    lang := Lang{"Go", 2009, "http://golang.org/"}

    // Reflection. %v => value
    fmt.Printf("%v\n", lang) // or %+v
    fmt.Printf("%+v\n", lang) // or %+v
    fmt.Printf("%#v\n", lang) // or %+v

    // Marshal into JSON
    data, err := json.Marshal(lang)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s\n", data);
}
