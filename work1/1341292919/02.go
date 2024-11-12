package main

import "fmt"

func main() {
	var height [10]int
	for i := 0; i < 10; i++ {
		fmt.Scan(&height[i])
	}
	var MaxHeight, number int
	fmt.Scan(&MaxHeight)
	for i := 0; i < 10; i++ {
		if MaxHeight+30 >= height[i] {
			number++
		}
	}
	fmt.Println(number)
}
