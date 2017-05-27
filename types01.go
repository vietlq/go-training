package main

import "fmt"

type Office int

const (
    Boston Office = iota
    NewYork
)

var officePlace = [...]string{
    "Boston, MA",
    "New York, NY",
}

func (o Office) String() string {
    return "Code2Pro, " + officePlace[o]
}

func main() {
    fmt.Printf("Hello, %s\n", Boston);
}
