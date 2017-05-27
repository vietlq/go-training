package main

import "fmt"

type Word struct {}

func (w *Word) String() string {
    return "Мир!";
}

func main() {
    fmt.Printf("Hello, %s\n", new(Word))
    //fmt.Printf("Hello, %s\n", "Мир!")
}
