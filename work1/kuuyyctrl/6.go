package main

import "fmt"

func main() {
	for i := 1; i <= 10; i++ {
		for j := i; j < 10; j++ {
			fmt.Printf("%d*%d =%2d ", i, j, (i * j))
		}
		fmt.Println()
	}
}
