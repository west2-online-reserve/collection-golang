package main

import "fmt"

func main() {
	var a [10]int
	var h, n int
	n = 0
	for i := 0; i < 10; i++ {
		fmt.Scan(&a[i])
	}
	fmt.Scan(&h)
	for i := 0; i < 10; i++ {
		if h+30 >= a[i] {
			n++
		}
	}
	fmt.Println(n)
}
