package main

import "fmt"

func main() {
    defer fmt.Println("world")
    defer fmt.Println("yahoo")

    fmt.Println("hello")

    defer fmt.Println("deferrals are stacked")
    defer fmt.Println("This the last but printed the first")
}
