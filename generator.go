package main

import "fmt"

func Count(start, end int) chan int {
	c := make(chan int)

	go func(ch chan int) {
		for i := start; i < end; i++ {
			ch <- i
		}

		close(ch)
	}(c)

	return c
}

// My take on this pattern is
// Sequence and blocking process,
// Useful if you need to build something sequential.
// And if there is the need to finish any prequisite process
func main() {

	fmt.Println("Here we go")
	for i := range Count(1, 10) {
		fmt.Printf("> %v \n", i)
		fmt.Println("one at the time ok")
	}

	fmt.Println("Finish")
}
