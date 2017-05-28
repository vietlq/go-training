package main

import (
    "fmt"
    "math"
)

func pow(x, n, lim float64) float64 {
    // Short statement before the condition, like for
    if v := math.Pow(x, n); v < lim {
        return v
    } else {
        fmt.Printf("%g >= %g\n", v, lim)
    }
    // v is not available after the if-else block
    return lim
}

func main() {
    fmt.Println(
        pow(3, 2, 10),
        pow(3, 3, 20),
    )
}
