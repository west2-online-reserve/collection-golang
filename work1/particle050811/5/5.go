package main

import "fmt"

func main() {
	a := make([]int, 0, 50)
	for i := 1; i <= 50; i++ {
		a = append(a, i)
	}
	for i := 0; i < len(a); {
		if a[i]%3 == 0 {
			a = append(a[:i], a[i+1:]...)
		} else {
			i++
		}
	}
	a = append(a, 114514)
	fmt.Println(a)
}
