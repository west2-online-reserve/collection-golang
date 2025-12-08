package main

import "fmt"

func main() {
	var slice []int = make([]int, 50)
	for i := 1; i <= 50; i++ {
		slice[i-1] = i
	}
	for i := 2; i < len(slice); i += 2 {
		slice = append(slice[0:i], slice[i+1:]...)
	}
	slice = append(slice, 114514)
	fmt.Println(slice)
}
