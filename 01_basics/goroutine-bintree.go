package main

import (
    "golang.org/x/tour/tree"
    "fmt"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
    if t.Left != nil {
        Walk(t.Left, ch)
    }
    ch <- t.Value
    if t.Right != nil {
        Walk(t.Right , ch)
    }
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
    c1 := make(chan int, 10)
    c2 := make(chan int, 10)
    b1, b2 := []int{}, []int{}
    x, y := 0, 0

    go Walk(t1, c1)
    go Walk(t2, c2)

    for {
        select {
        case x = <-c1:
            b1 = append(b1, x)
        case y = <-c2:
            b2 = append(b2, y)
        }
        if len(b1) >= 10 && len(b2) >= 10 {
            break
        }
    }

    if len(b1) != len(b2) {
        return false
    }

    for i := range b1 {
        if b1[i] != b2[i] {
            return false
        }
    }
    return true
}

func main() {
    t1, t2, t3 := tree.New(1), tree.New(1), tree.New(2)
    fmt.Println("t1 == t2:", true == Same(t1, t2))
    fmt.Println("t1 != t3:", false == Same(t1, t3))
}
