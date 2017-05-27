package main

import (
    "fmt"
    "encoding/json"
    "log"
    "os"
    "io"
)

type Lang struct {
    Name string
    Year int
    URL string
}

func main() {
    input, err := os.Open("lang.json")
    if err != nil {
        log.Fatal(err)
    }

    dec := json.NewDecoder(input)
    for {
        var lang Lang
        err := dec.Decode(&lang)
        if err != nil {
            if err == io.EOF {
                break
            }
            log.Fatal(err)
        }
        fmt.Printf("%v\n", lang)
    }
}
