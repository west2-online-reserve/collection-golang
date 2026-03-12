package main

import "fmt"

func main(){

	a := make([]int, 51)
	var b []int

	for i := 1; i <= 50; i++ {
		a[i] = i
	}

	for _, v := range a {
		if v % 3 == 0 {
			continue
		}
		b = append(b, v)
	}
	b = append(b, 114514)

	for _, v := range b {
		fmt.Printf("%d ", v)
	}
}