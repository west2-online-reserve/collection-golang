package main

import "fmt"

func main() {
	var slice []int
	for i := 1; i <= 50; i++ {
		slice = append(slice, i)
	}
	for i := 3; i < len(slice); i += 2 {
		slice = append(slice[:i-1], slice[i:]...)
	}
	slice = append(slice, 114514)
	fmt.Println(slice)
}
