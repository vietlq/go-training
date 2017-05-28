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

// Declare methods on non-struct types
type MyFloat float64
func (f MyFloat) Abs() float64 {
    if f < 0 {
        return float64(-f)
    }
    return float64(f)
}
// You can only declare a method with a receiver whose type is defined in the same package as the method. You cannot declare a method with a receiver whose type is defined in another package (which includes the built-in types such as int).

// Method with pointer receiver
func (v *Vertex) Scale(f float64) {
    v.X *= f
    v.Y *= f
}

func main() {
    v := Vertex{3, 4}
    fmt.Println(v.Abs())
    v.Scale(3)
    fmt.Println(v)

    f := MyFloat(-1.2345)
    fmt.Println(f.Abs())
}
