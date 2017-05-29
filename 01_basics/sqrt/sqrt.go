package main

import (
    "fmt"
    "math"
)

func Abs(x float64) float64 {
    if x < 0 {
        return -x
    }
    return x
}

func Sqrt(x float64) float64 {
    if x < 0 {
        return math.NaN()
    }
    epsilon := 1e-6
    z := x
    count := 0
    for epsilon < Abs(z - x/z) {
        count++
        z = z - (z*z - x)/(2*z)
    }
    fmt.Printf("It took me %d loops\n", count)
    return z
}

func main() {
    fmt.Println("Sqrt(2) =", Sqrt(2))
}
