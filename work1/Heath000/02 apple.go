package main

import (
	"fmt"
)

func main() {
	var (
		a      [10]int
		height int
	)
	for i := 0; i < 10; i++ {
		fmt.Scan(&a[i])
	}
	fmt.Scan(&height)
	n := 0
	height += 30
	for i := 0; i < 10; i++ {
		if a[i] <= height {
			n++
		}
	}
	fmt.Println(n)
}
