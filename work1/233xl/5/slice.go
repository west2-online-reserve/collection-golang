package main

import "fmt"

func main() {
	oneToFifty := make([]int, 0, 50)
	for i := 1; i <= 50; i++ {
		if i % 3 != 0 {
			oneToFifty = append(oneToFifty, i)
		}
	}
	oneToFifty = append(oneToFifty, 114514)

	fmt.Print(oneToFifty)
}