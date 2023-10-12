package main

import "fmt"

func main() {
	var list [10]int
	for i := 0; i <= 9; i++ {
		fmt.Scanf("%d", &list[i])
	}
	var h, a int
	fmt.Scan(&h)
	h += 30
	for i := 0; i <= 9; i++ {
		if list[i] <= h {
			a++
		}
	}
	fmt.Println(a)
}
