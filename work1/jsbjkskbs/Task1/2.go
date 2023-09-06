package main

import (
	"fmt"
)

func main() {
	var apples [10]int
	var reachable int
	for i := 0; i < 10; i++ {
		fmt.Scan(&apples[i])
	}
	fmt.Scan(&reachable)
	reachable = reachable + 30
	var cnt = 0
	for i := 0; i < 10; i++ {
		if apples[i] <= reachable {
			cnt = cnt + 1
		}
	}
	fmt.Print(cnt)
}
