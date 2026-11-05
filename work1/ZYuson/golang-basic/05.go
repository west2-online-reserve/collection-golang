package main

import "fmt"

func main() {
	var slice []int
	for i := 0; i <= 50; i++ {
		slice = append(slice, i)
	}
	var NewSlice []int
	for _, v := range slice {
		if v%3 != 0 {
			NewSlice = append(NewSlice, v)
		}
	}
	NewSlice = append(NewSlice, 114514)
	fmt.Println(NewSlice)
}
