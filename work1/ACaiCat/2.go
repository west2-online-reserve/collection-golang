package main

import "fmt"

func main() {
	const (
		chairHeight = 30
		appleCount  = 10
	)

	var appleHeights [appleCount]int
	var taoHeight int

	for i := 0; i < appleCount; i++ {
		_, _ = fmt.Scan(&appleHeights[i])
	}

	_, _ = fmt.Scan(&taoHeight)

	var reached int
	for _, appleHeight := range appleHeights {
		if appleHeight <= taoHeight+chairHeight {
			reached++
		}
	}
	fmt.Println(reached)

}
