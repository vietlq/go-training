package main

import "fmt"
import "time"

/*
Go programs express error state with error values.
The error type is a built-in interface similar to fmt.Stringer:
type error interface {
    Error() string
}
*/

type MyError struct {
    When time.Time
    What string
}

// Implement the method Error() string
func (e *MyError) Error() string {
    return fmt.Sprintf("at %v, %s", e.When, e.What)
}

func run() error {
    return &MyError{
        time.Now(),
        "it didn't work",
    }
}

func main() {
    if err := run(); err != nil {
        fmt.Println(err)
    }
}
