package main

import "fmt"

func sum(s []int, c chan int) {
    sum := 0
    for _, v := range s {
        sum += v
    }
    c <- sum // send sum to c (yield)
}

func main() {
    s := []int{10, 9, 13, 15, 20}

    // Channel of ints. Unbuffered
    c := make(chan int)
    go sum(s[:len(s)/2], c)
    go sum(s[len(s)/2:], c)
    x, y := <-c, <-c // Receive from c

    fmt.Println(x, y, x + y)

    // Buffered channel with 100 slots
    ch := make(chan int, 100)
    ch <- 1
    ch <- 2
    fmt.Println(<-ch)
    fmt.Println(<-ch)
}
