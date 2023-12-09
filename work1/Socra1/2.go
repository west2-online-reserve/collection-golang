package main

import (
	"fmt"
)

func main() {
	var a [10]int
	for i := 0; i < 10; i++ {
		fmt.Scan(&a[i])
	}
	var hgt int
	cnt := 0
	fmt.Scan(&hgt)
	hgt += 30
	for i := 0; i < 10; i++ {
		if hgt >= a[i] {
			cnt++
		}
	}
	fmt.Printf("%v\n", cnt)
}
