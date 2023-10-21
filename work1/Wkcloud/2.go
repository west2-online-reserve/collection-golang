package main

import "fmt"

func main() {
	height := [10]int{}
	var reach, x int
	for i := 0; i < 10; i++ {
		fmt.Scanf("%d", &height[i])
	}
	fmt.Scan(&reach)
	reach += 30
	for _, n := range height {
		if n <= reach {
			x++
		}
	}
	fmt.Printf("%d", x)
}
