package main

import (
    "fmt"
    "time"
)

func fetch(url string) {
}

func main() {
    start := time.Now()
    fetch("http://code2.pro")
    elapsed := time.Since(start)
    fmt.Println(elapsed)
    fmt.Printf("%d ns\n", elapsed)
}
