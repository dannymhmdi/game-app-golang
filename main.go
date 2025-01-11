package main

import (
	"fmt"
	"regexp"
)

func main() {
	sl := [4]int{1, 2, 3, 4}
	nl := sl[:]
	nl[0] = 10
	fmt.Println(sl, nl)
	str := "091272752362"
	match, _ := regexp.MatchString(`^0\d{10}$`, str)
	fmt.Println(match)
	fmt.Printf("%T\n", str[2:])
}
