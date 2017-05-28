package main

import "fmt"

func main() {
    primes := [7]int{2, 3, 5, 7, 11, 13, 17}

    // Slice: []T
    var s []int = primes[1:4]
    fmt.Println(primes)
    fmt.Println(s)

    names := [4]string{
        "John",
        "Misha",
        "Alice",
        "Claire",
    }
    fmt.Println(names)

    // Slices do not store any data, they are merely views/references with boundaries
    a := names[0:2]
    b := names[1:3]
    b[0] = "Misha the Great"
    fmt.Println(a, b)
    fmt.Println(names)
}
