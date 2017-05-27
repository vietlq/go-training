package main

import (
    "encoding/hex"
    "fmt"
    "os"
)

func main() {
    h := hex.Dumper(os.Stdout)
    defer h.Close()
    fmt.Fprintf(h, "Hello, World!\n")
}
