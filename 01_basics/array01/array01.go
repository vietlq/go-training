package main

import "fmt"

type Location int

const (
    London Location = iota
    NewYork
    Singapore
    Frankfurt
)

// Unbounded array. Feel free to add trailing comma
var locationMap = [...]string{
    "London, UK",
    "New York, US",
    "Singapore, SG",
    "Frankfurt, DE",
}

// Unnamed function that returns String() string
func (l Location) String() string {
    return fmt.Sprintf("Location: %s", locationMap[l])
}

func main() {
    var a [2]string
    a[0] = "Hello"
    a[1] = "World"
    fmt.Println(a[0], a[1])
    fmt.Println(a)

    primes := [7]int{2, 3, 5, 7, 11, 13, 17}
    fmt.Println(primes)

    fmt.Printf("Type: %T\nValue: %v\nStr: %s\n", London, London, London)
}
