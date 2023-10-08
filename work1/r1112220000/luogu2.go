package main

import "fmt"

func main() {
	var arr = make([]int, 10)
	var high int
	var n int = 0

	for i := 0; i < 10; i++ {
		_, _ = fmt.Scan(&arr[i])
	}
	_, _ = fmt.Scan(&high)
	for i := 0; i < 10; i++ {
		if arr[i] <= high+30 {
			n++
		}
	}
	fmt.Printf("%d", n)
}
