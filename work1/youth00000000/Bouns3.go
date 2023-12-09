package main

import "fmt"

func main() {
	num := [...]int{3, 2, 4}
	var b []int
	a := num[:]
	var target = 6
	for i := 0; i < len(a); i++ {
		for j := i + 1; j < len(a); j++ {
			if a[i]+a[j] == target {
				b = append(b, i, j)
				fmt.Println(b)
			}
		}
	}
}
