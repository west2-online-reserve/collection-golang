package main

import (
	"fmt"
)

func main() {
	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			fmt.Printf("%d%s%d%s%d ", j, "X", i, "=", i*j)
		}
		fmt.Println()
	}
}
