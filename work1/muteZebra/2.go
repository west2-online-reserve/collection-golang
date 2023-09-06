package main

import "fmt"

func main() {
	var height, total int
	apples := make([]int, 10)

	for i := 0; i < 10; i++ {
		_, _ = fmt.Scan(&apples[i])
	}

	_, _ = fmt.Scan(&height)

	for i := 0; i < 10; i++ {
		if height+30 >= apples[i] {
			total++
		}
	}
	fmt.Println(total)
}
