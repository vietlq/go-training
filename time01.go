package main

import (
    "fmt"
    "time"
)

func main() {
    day := time.Now().Weekday()
    fmt.Printf("Hello, %s (%d)\n", day, day)
}
