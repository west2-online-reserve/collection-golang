package main

import "fmt"

func main() {
	var apples [10]int
	for i := 0; i < 10; i++ {
		fmt.Scan(&apples[i])
	}
	var height int
	fmt.Scan(&height)
	height += 30
	count := 0
	for i := 0; i < 10; i++ {
		if apples[i] <= height {
			count++
		}
	}
	fmt.Println(count)
}
