package main

import "fmt"

// Variables at package level
var c, python, java bool

func main() {
    // Variable at function level
    var i int
    // See the initialised values: 0 & false
    fmt.Println(i, c, python, java)
}
