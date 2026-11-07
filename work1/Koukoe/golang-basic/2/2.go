package main

import "fmt"

func main() {
	var aHeights [10]int
	var tHeight, count int
	var bench int = 30

	for i := range aHeights {
		fmt.Scan(&aHeights[i])
	}

	fmt.Scan(&tHeight)
	for _, aHeight := range aHeights {
		if aHeight <= tHeight+bench {
			count++
		}
	}

	fmt.Println(count)
}
