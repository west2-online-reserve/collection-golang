package main

import "fmt"

func main() {
	var height [10]int
	var i, j, k int
	for i = 0; i < 10; i++ {
		fmt.Scanf("%d", &j)
		height[i] = j
	}
	fmt.Scanf("%d\n", &k)
	var cnt = 0
	for i = 0; i < 10; i++ {
		if height[i] < (j + 30) {
			cnt++
		}

	}
	fmt.Println(cnt)
}
