package main

import "fmt"

func add (x, y int) int {
    return x + y
}

func main() {
    fmt.Printf("%d + %d = %d\n", 2, 3, add(2, 3))
    fmt.Println(4, "+", 5, "=", add(4, 5))
}
