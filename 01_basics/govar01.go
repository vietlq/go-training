package main

import "fmt"

// Variables at package level
var c, python, java bool

func main() {
    // Variable at function level
    var i int
    // See the initialised values: 0 & false
    fmt.Println(i, c, python, java)
    // Variables with initialisers
    var t, v int = 1, 2
    var x, y, z = true, false, "yes!"
    fmt.Println(t, v, x, y, z)
}
