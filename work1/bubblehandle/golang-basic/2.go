package main

import "fmt"

func main() {
	var arr [10]int
	var height int
	for i := 0; i < 10; i++ {
		fmt.Scan(&arr[i])
	}
	fmt.Scan(&height)
	height += 30
	var ans int = 0
	for i := 0; i < 10; i++ {
		if height >= arr[i] {
			ans++
		}
	}
	fmt.Println(ans)
}
