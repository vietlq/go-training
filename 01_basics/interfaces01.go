package main

import (
    "fmt"
    "math"
)

// An interface type is defined as a set of method signatures.
// A value of interface type can hold any value that implements those methods.
type Abser interface {
    Abs() float64
}
// Interfaces are implemented implicitly
// A type implements an interface by implementing its methods. There is no explicit declaration of intent, no "implements" keyword.
// Implicit interfaces decouple the definition of an interface from its implementation, which could then appear in any package without prearrangement.

func main() {
    // Interface value. Equivalent to C++ pure virtual base class instance
    var a Abser

    f := MyFloat(-math.Sqrt2)
    v := Vertex{3, 4}

    a = f
    fmt.Println(a.Abs())
    describe(a)

    a = &v
    fmt.Println(a.Abs())
    describe(a)
    // This will fail. Implement func (v Vertex) Abs() float64
    //a = v
}

type MyFloat float64

func (f MyFloat) Abs() float64 {
    if f < 0 {
        return float64(-f)
    }
    return float64(f)
}

type Vertex struct {
    X, Y float64
}

func (v *Vertex) Abs() float64 {
    return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Pass interface value into the function
func describe(a Abser) {
    fmt.Printf("(%v, %T)\n", a, a)
}
