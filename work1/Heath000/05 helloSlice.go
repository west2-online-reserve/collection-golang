package main

import (
	"fmt"
)

func main() {
	a := make([]int, 50, 50)
	for i := 0; i < 50; i++ {
		a[i] = i + 1
	}
	for i := 0; i < len(a); i++ {
		if a[i]%3 == 0 {
			a = append(a[:i], a[i+1:]...)
		}
	}
	a = append(a, 114514)
	for i := 0; i < len(a); i++ {
		fmt.Println(a[i])
	}
}
