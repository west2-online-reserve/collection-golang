package main

import "fmt"

func main() {
	var arr [10]int
	var n int
	num := 0

	for i := 0; i < 10; i++ {
		fmt.Scanf("%d", &arr[i])
	}
	fmt.Scanf("\n%d", &n)
	n = n + 30
	for j := 0; j < 10; j++ {
		if arr[j] <= n {
			num += 1
		}
	}
	fmt.Printf("%d\n", num)
}
