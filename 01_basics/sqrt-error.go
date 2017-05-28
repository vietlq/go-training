package main

import "fmt"
import "math"

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
    // Note: a call to fmt.Sprint(e) inside the Error method will send the program into an infinite loop. You can avoid this by converting e first: fmt.Sprint(float64(e)). Why?
    return fmt.Sprintf("cannot Sqrt negative number: %f", float64(e))
}

func Sqrt(x float64) (float64, error) {
    if x < 0 {
        return -1, ErrNegativeSqrt(x)
    }
    return math.Sqrt(x), nil
}

func main() {
    fmt.Println(Sqrt(-2))
    fmt.Println(Sqrt(2))
}
