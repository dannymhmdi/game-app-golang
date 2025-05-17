package main

import (
	"fmt"
	"math"
)

func main() {
	//param := []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	//fmt.Printf("%+v\n", batchGenerator(2, param))
	testing("Danial", "Shirin", "Sadra")
}

func batchGenerator(size int, slc []uint) [][]uint {
	batchCount := int(math.Ceil(float64(len(slc)) / float64(size)))
	//remain := len(slc) % size
	ultSlc := make([][]uint, batchCount)

	for i := 0; i < batchCount; i++ {
		start := i * size
		end := start + size
		fmt.Println("start", start, "end", end)
		if end > len(slc) {
			end = len(slc)
		}

		ultSlc[i] = slc[start:end]
	}

	return ultSlc
}

func testing(fields ...string) {
	fmt.Printf("fields:%+v,type:%T\n", fields, fields)
}
