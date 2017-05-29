package main

import (
    "fmt"
    "time"
)

func main() {
    t := time.Now()
	// Better way to write long & nested if-else blocks
    switch {
    case t.Hour() < 12:
        fmt.Println("Good morning!")
    case t.Hour() < 17:
        fmt.Println("Good afternoon.")
    default:
        fmt.Println("Good evening.")
    }
}
