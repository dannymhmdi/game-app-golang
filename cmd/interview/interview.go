package main

import "fmt"

func main() {
	slc := []int{10, 15, 3, 7}
	k := 25
	s := Finder(slc, k)
	fmt.Println("isExist", s)

}

func Finder(slc []int, k int) int {
	for _, num := range slc {
		for _, numInner := range slc {
			if num == numInner {
				continue
			}
			sum := num + numInner
			if sum == k {
				return sum
			}
		}
	}
	return -1
}
