package main

import "fmt"

func main() {
	var a []int
	for i := 1; i <= 50; i++ {
		a = append(a, i)
	}
	for i, j := range a {
		if j%3 == 0 {
			a = append(a[:i], a[i+1:]...)
		}
	}
	a = append(a, 114514)
	fmt.Print(a)
}
