package main

import (
	"fmt"
	"math"
)

func main() {
	//param := []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	//fmt.Printf("%+v\n", batchGenerator(2, param))
	//testing("Danial", "Shirin", "Sadra")
	//str := "refresh-token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjIsIk5hbWUiOiJEYW5pZWwiLCJSb2xlIjoyLCJSZWdpc3RlcmVkQ2xhaW1zIjp7InN1YiI6InJ0IiwiZXhwIjoxNzQ5MTA5MTgxfX0.1Q8uY0-b-PAqAZ3l-U2TQnrLaqwkCYPCU86IvOqDtlg"
	//refershToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjIsIk5hbWUiOiJEYW5pZWwiLCJSb2xlIjoyLCJSZWdpc3RlcmVkQ2xhaW1zIjp7InN1YiI6InJ0IiwiZXhwIjoxNzQ5NjY0NDIyfX0.qj_NS7j8tXislIMNlRrF3eC4S04dwFREv9uLyHwgYps"
	hashedToken := "$2a$10$QEpoectYwcyl1nEI2QpSvu/btIljfAiW9.2r6iuXH4A9Nw5MGmp3i"
	tokenDB := "$2a$10$6UaCfti4wq5sUmBpVUUz..UO4XNbaYMcaUCR6mUcUBvbnwG0aGHd2"
	fmt.Println("compare", hashedToken == tokenDB)
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
