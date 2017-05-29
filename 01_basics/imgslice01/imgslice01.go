package main

// Run the command: go get golang.org/x/tour/gotour
import "golang.org/x/tour/pic"

func Pic(dx, dy int) [][]uint8 {
    img := make([][]uint8, dy)
    for i := range img {
        img[i] = make([]uint8, dx)
    }
    return img
}

func main() {
    pic.Show(Pic)
}
