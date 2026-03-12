package main

import (
	"fmt"
)

func main() {
	var slice = make([]int, 50)
	for i := 0; i < 50; i++ {
		slice[i] = i + 1
	}
	var newSlice = make([]int, 0)
	for i := 0; i < 50; i++ {
		if (i+1)%3 == 0 {
			continue
		}
		newSlice = append(newSlice, slice[i])
	}
	newSlice = append(newSlice, 114514)
	fmt.Println(newSlice)
}
