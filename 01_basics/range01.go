package main

import "fmt"

var pow = []int{1, 2, 4, 8, 16, 32}

func main() {
    for i, v := range pow {
        fmt.Printf("2**%d = %d\n", i, v)
    }

    for i := range pow {
        fmt.Println("i =", i)
    }

    for _, v := range pow {
        fmt.Println("v = ", v)
    }
}
