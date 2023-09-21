package main

import "fmt"

func main() {
	var height [10]int
	for i := 0; i < 10; i++ {
		_, _ = fmt.Scan(&height[i])
	}
	var h int
	_, _ = fmt.Scan(&h)
	ans := 0
	for i := 0; i < 10; i++ {
		if h+30 >= height[i] {
			ans++
		}
	}
	fmt.Println(ans)
}
