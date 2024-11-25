package main

import (
	"fmt"
)

func main() {
	A := make([]int, 10)
	for i := 0; i < 10; i++ {
		fmt.Scan(&A[i])
	}
	var B int
	fmt.Scan(&B)
	B += 30
	var tag int
	for _, value := range A {
		if B >= value {
			tag++
		}
	}
	fmt.Println(tag)
}
