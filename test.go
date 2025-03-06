package main

import "fmt"

func main() {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	close(ch)

	// Receiving from the closed channel
	for value := range ch {
		fmt.Println("value", value)
	}
	fmt.Println("done")
}
