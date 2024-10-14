package main

import (
	"fmt"
)

func main() {
	var height [10]int
	var H, s int

	for i := 0; i < 10; i++ {
		fmt.Scan(&height[i])
	}

	fmt.Scan(&H)
	H += 30

	for i := 0; i < 10; i++ {
		if H >= height[i] {
			s++
		}
	}

	fmt.Println(s)
}
