package main

import "fmt"

func main() {
	var height [10]int
	for i := 0; i < 10; i++ {
		fmt.Scan(&height[i])
	}
	var canget int
	var num int = 0
	fmt.Scan(&canget)
	canget = canget + 30
	for i := 0; i < 10; i++ {
		if height[i] <= canget {
			num++
		}
	}
	fmt.Println(num)
}
