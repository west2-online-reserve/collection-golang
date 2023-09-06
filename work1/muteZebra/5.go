package main

import (
	"fmt"
	"math/rand"
)

func main() {
	a := make([]int, 50)
	var result []int
	for i := 0; i < 50; i++ {
		a[i] = rand.Intn(50)
	}

	for _, v := range a {
		if v%3 != 0 {
			result = append(result, v)
		}
	}
	result = append(result, 114514)
	for _, v := range result {
		fmt.Printf("%d  ", v)
	}
}
