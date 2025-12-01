package main

import "fmt"

func main() {
	var k int
	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			k = i * j
			fmt.Printf("%d*%d=%-4d", j, i, k)
		}
		fmt.Printf("\n")
	}
}
