package main

import "fmt"

func main() {
	var num [10]int
	for i := 0; i < 10; i++ {
		fmt.Scanf("%d", &num[i])
	}
	var height int
	fmt.Scanf("%d", &height)
	cnt := 0
	for i := 0; i < 10; i++ {
		if (height + 30) >= num[i] {
			cnt++
		}
	}
	fmt.Printf("%d", cnt)
}