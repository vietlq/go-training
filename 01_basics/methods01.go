package main

import (
    "fmt"
    "math"
)

type Vertex struct {
    X, Y float64
}

/*
Go does not have classes. However, you can define methods on types.
A method is a function with a special receiver argument.
The receiver appears in its own argument list between the func keyword and the method name.
In this example, the Abs method has a receiver of type Vertex named v.
*/

// Instances of the struct Vertex now have a method called .Abs()
func (v Vertex) Abs() float64 {
    return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func main() {
    v := Vertex{3, 4}
    fmt.Println(v.Abs())
}
