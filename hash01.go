package main

import (
    "fmt"
    "hash/crc32"
)

func main() {
    // h is a writer! You can write to it!
    h := crc32.NewIEEE()
    fmt.Fprintf(h, "Hello, World!")
    fmt.Printf("hash=%#x\n", h.Sum32())
}
