package main

import "fmt"

type Vertex struct {
    X int
    Y int
}

// Struct literals
var (
    v1 = Vertex{1, 2}  // has type Vertex
    v2 = Vertex{X: 1}  // Y: 0 is implicit
    v3 = Vertex{}      // implicit X:0 and Y:0
    q  = &Vertex{3, 4} // has type *Vertex
)

func main() {
    v := Vertex{1, 2}
    fmt.Println(v)
    v.X = 4
    fmt.Println(v)
    p := &v
    p.Y = 8
    fmt.Println(v)
}
