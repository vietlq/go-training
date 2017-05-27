package main

import (
    "fmt"
    "encoding/json"
    "log"
    "os"
    "io"
    "io/ioutil"
    "time"
    "net/http"
)

type Lang struct {
    Name string
    Year int
    URL string
}

// Higher-order function
func do(f func(Lang)) {
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
        f(lang)
    }
}

func count(name, url string) {
    start := time.Now()
    r, err := http.Get(url)
    if err != nil {
        fmt.Printf("%s: %s", name, err)
        return
    }
    n, _ := io.Copy(ioutil.Discard, r.Body)
    r.Body.Close()
    fmt.Printf("%s %d [%.2fs]\n", name, n, time.Since(start).Seconds())
}

func main() {
    // Pass function literal
    start := time.Now()
    do(func(lang Lang) {
        count(lang.Name, lang.URL)
    })
    fmt.Printf("%.2fs total\n", time.Since(start).Seconds())
}
