package main

import "fmt"

func main() {
	var slice []int

	for i := 1; i <= 50; i++ {
		if i%3 == 0 {
			continue
		}
		slice = append(slice, i)
	}
	slice = append(slice, 114514)

	fmt.Println(slice)
}
