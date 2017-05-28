package main

import "fmt"

func main() {
    sum := 0
    for i := 0; i < 10; i++ {
        sum += i
    }
    fmt.Println(sum)

    sum = 1
    // Think of it as a while loop
    // Init & post statements are optional
    for ; sum < 1000; {
        sum += sum
    }
    fmt.Println(sum)

    // Now drop the semicolons. There's no while in Go
    sum = 1
    for sum < 100 {
        sum += sum
    }
    fmt.Println(sum)
}
