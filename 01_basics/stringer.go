package main

import "fmt"

type Person struct {
    Name string
    Age int
}

// Implement Stringer interface by defining the String() string method
func (p Person) String() string  {
    return fmt.Sprintf("Name: %s, Age: %d", p.Name, p.Age)
}

func main()  {
    a := Person{"Misha", 14}
    b := Person{"Claire", 20}
    fmt.Println(a)
    fmt.Println(b)
}