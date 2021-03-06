package main

import "fmt"

// This exercise shows that s holds ref to an array and you can adjust view dynamically
func main() {
    s := []int{2, 3, 5, 7, 11, 13, 17}
    printSlice(s)

    // Slice the slice to give it zero length
    s = s[:0]
    printSlice(s)

    // Extend its length
    s = s[:4]
    printSlice(s)

    // Drop first 2 values
    s = s[2:]
    printSlice(s)

    var v []int
    fmt.Println(v, len(v), cap(v))
    if v == nil {
        fmt.Println("Nil slice!")
    }
}

func printSlice(s []int) {
    fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}
