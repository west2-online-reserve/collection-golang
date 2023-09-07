package main

import "fmt"

func main() {
	var list [10]int
	for i := 0; i <= 9; i++ {
		_, _ = fmt.Scan(&list[i])
	}

	var height int
	_, _ = fmt.Scan(&height)
	height += 30

	count := 0
	for i := 0; i <= 9; i++ {
		if list[i] <= height {
			count++
		}
	}

	fmt.Println(count)
}
