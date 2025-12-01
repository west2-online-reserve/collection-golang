package main

import "fmt"

func main() {
	var treeHeights [10]int
	var taoHeight int
	for i := 0; i < 10; i++ {
		fmt.Scan(&treeHeights[i])
	}
	fmt.Scan(&taoHeight)
	taoHeight += 30
	count := 0
	for i := 0; i < 10; i++ {
		if treeHeights[i] <= taoHeight {
			count++
		}
	}
	fmt.Printf("%d", count)
}
