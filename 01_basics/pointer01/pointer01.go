package main

import "fmt"

func main() {
    var p *int
    i := 20
    p = &i
    fmt.Println("i =", *p);
    q := &i
    *q *= 21
    fmt.Println("i =", *p);
}
