package main

import (
	"fmt"
	"sync"
)

func Merge(clients ...<-chan int) <-chan int {
	var wg sync.WaitGroup

	out := make(chan int)

	send := func(c <-chan int) {
		for n := range c {
			out <- n
		}

		wg.Done()
	}

	fmt.Printf("Create %v channels", len(clients))
	wg.Add(len(clients))

	for _, c := range clients {
		go send(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {

	c1 := make(<-chan int)
	c2 := make(<-chan int)

	out := Merge(c1, c2)

	for o := range out {
		fmt.Println(o)
	}
}
