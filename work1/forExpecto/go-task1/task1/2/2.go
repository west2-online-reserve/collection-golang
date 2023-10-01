package main

import "fmt"

func main() {
	var a [10]int
	var maxheight int
	num := 0
	for i := 0; i < 10; i++ {
		fmt.Scan(&a[i])
	}
	fmt.Scan(&maxheight)
	for i := 0; i < 10; i++ {
		if maxheight+30 >= a[i] {
			num++
		}
	}
	fmt.Println(num)
}
