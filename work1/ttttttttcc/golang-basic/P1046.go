package main

import (
	"fmt"
)

func main() {
	var appleheights []int
	var maxStretch, height int

	for i := 0; i < 10; i++ {
		fmt.Scan(&height)
		appleheights = append(appleheights, height)
	}

	fmt.Scan(&maxStretch)

	maxReach := maxStretch + 30

	count := 0
	for _, h := range appleheights {
		if h <= maxReach {
			count++
		}
	}

	fmt.Println(count)
}