package main

import "fmt"

func main() {
	var heights [10]int
	for i := range heights {
		fmt.Scan(&heights[i])
	}
	var reachHeight int
	fmt.Scan(&reachHeight)
	reachHeight += 30
	count := 0
	for _, height := range heights {
		if height <= reachHeight {
			count++
		}
	}
	fmt.Println(count)
}
