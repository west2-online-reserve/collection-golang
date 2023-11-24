package main

import "fmt"

func main() {
	var slice []int
	for i := 0; i < 50; i++ {
		slice = append(slice, i+1)

	}
	for i, v := range slice {
		if v%3 == 0 {
			slice = append(slice[:i], slice[i+1:]...)
		}
	}
	slice = append(slice, 114514)
	fmt.Printf("slice: %v\n", slice)
}
