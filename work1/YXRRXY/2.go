package main

import "fmt"

func main() {
	var heights [11]int
	var reach_h, ans int
	for i:=1; i<=10; i++ {
		fmt.Scan(&heights[i])
	}
	fmt.Scan(&reach_h)
	reach_h += 30
	for i:=1; i<=10; i++ {
		if heights[i] <= reach_h {
			ans++
		}
	}
	fmt.Println(ans)
}