package main

import "fmt"

func main() {

	var heights [10]int
	var maxReach int
	var counts int

	for i := 0; i <= 9; i++ {
		fmt.Scan(&heights[i])
	}

	fmt.Scan(&maxReach)
	maxReach += 30

	for _, i := range heights {
		if maxReach >= i {
			counts++
		}

	}
	fmt.Println(counts)

}
