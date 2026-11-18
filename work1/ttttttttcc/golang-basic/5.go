package main

import (
	"fmt"
)

func main() {
	slice := make([]int, 0)
	for i := 1; i <= 50; i++ {
		slice = append(slice, i)
	}

	var newSlice []int
	for _, num := range slice {
		if num%3 != 0 {
			newSlice = append(newSlice, num)
		}
	}

	newSlice = append(newSlice, 114514)

	fmt.Println(newSlice)
}