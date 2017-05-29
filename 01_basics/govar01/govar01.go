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
    // When using initiliser, the type can be omitted and inferred
    var x, y, z = true, false, "yes!"
    fmt.Println(t, v, x, y, z)
    // Short variable declaration only available inside functions
    a, b := 123, "good!"
    fmt.Println(a, b)
}
