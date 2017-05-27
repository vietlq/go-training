package main

import (
    "fmt"
    "encoding/xml"
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

    // Marshal into XML
    //data, err := xml.Marshal(lang)
    data, err := xml.MarshalIndent(lang, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s\n", data);
}
