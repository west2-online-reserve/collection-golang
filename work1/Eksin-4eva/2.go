package main

import "fmt"

func main() {
	var a [10]int
	var h, count int
	for i := 0; i < 10; i++ {
		fmt.Scan("%d", &a[i])
	}
	fmt.Scan("%d", h)
	for _, height := range a {
		if height < h+30 {
			count++
		}
	}
	fmt.Println(count)
}
