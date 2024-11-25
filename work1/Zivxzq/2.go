package main

import (
	"fmt"
)

func main() {
	var a int
	var num [11]int
	var res int
	for i := 1; i <= 10; i++ {
		fmt.Scan(&num[i])
	}
	fmt.Scan(&a)
	a += 30
	for i := 1; i <= 10; i++ {
		if num[i] <= a {
			res++
		}
	}
	fmt.Println(res)
}
