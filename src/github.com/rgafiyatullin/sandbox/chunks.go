package main

import "fmt"

func main() {
	chunk := []byte{0, 1, 2, 3, 4, 5}
	ch1 := chunk[0:0]
	ch2 := chunk[6:]

	fmt.Printf("%v %v\n", ch1, ch2)
}
