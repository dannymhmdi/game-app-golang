package main

import (
	"fmt"
	"time"
)

func main() {
	var u any
	fmt.Println("user", u)
	fmt.Println("time", time.Now().UnixMicro())
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	close(ch)
	s := []int{1, 2, 3, 4, 5}

	for i := 1; i < len(s); i = i + 2 {
		if i+1 <= len(s) {
			matchedUsers := struct {
				user1 uint
				user2 uint
			}{
				user1: uint(i),
				user2: uint(i + 1),
			}

			fmt.Printf("matchedUsers:%+v\n", matchedUsers)
		}

	}
	// Receiving from the closed channel
	for value := range ch {
		fmt.Println("value", value)
	}
	fmt.Println("done")
}

func Add(d time.Duration) time.Time {
	return time.Now().Add(d)
}
