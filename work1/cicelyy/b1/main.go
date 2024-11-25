package main

import "fmt"

func main() {
	const max = 9
	for i := 1; i <= max; i++ {
		for j := 1; j <= i; j++ {
			fmt.Printf("%d*%d=%d\t", j, i, i*j)
		}
		fmt.Println()
	}
}
