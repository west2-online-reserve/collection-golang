package main

import (
	"fmt"
	"sort"
)

func main() {
	arr := make([]int, 10)
	for i := 0; i < 10; i++ {
		fmt.Scan(&arr[i])
	}
	var x int
	fmt.Scan(&x)
	x += 30
	sort.Ints(arr)
	l, r := 0, 9
	for l < r {
		m := l + (r-l)/2
		if arr[m] > x {
			r = m
		} else {
			l++
		}
	}
	fmt.Println(l)
}
