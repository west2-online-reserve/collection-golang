package main

import "fmt"

func main() {
	var slice = make([]int, 10)
	for i := 0; i < 10; i++ {
		fmt.Scan(&slice[i])
	}
	var high int
	fmt.Scan(&high)
	high += 30
	count := 0
	for i := 0; i < 10; i++ {
		if high >= slice[i] {
			count++
		}
	}
	fmt.Println(count)
}
