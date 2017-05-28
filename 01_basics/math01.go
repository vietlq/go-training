package main

import (
    "math"
    "math/rand"
    "fmt"
    "time"
)

func main() {
    rand.Seed(time.Now().UTC().UnixNano())
    n := math.Abs(rand.NormFloat64()*10001)
    fmt.Printf("sqrt(%.2f) = %.4f\n", n, math.Sqrt(n))
}
