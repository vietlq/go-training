package main

import (
    "fmt"
    "hash/crc32"
    "io"
    "os"
)

func main() {
    // h is a writer! You can write to it!
    h := crc32.NewIEEE()
    // w is a multi-writer. Equivalent to Unix tee
    w := io.MultiWriter(h, os.Stdout)
    fmt.Fprintf(w, "Hello, World!\n")
    fmt.Printf("hash=%#x\n", h.Sum32())
}
