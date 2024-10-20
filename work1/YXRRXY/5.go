package main

import "fmt"

func main() {
	slice1 := make([]int, 0, 52)
	slice2 := make([]int, 0, 52)
	for i := 1; i <= 50; i++ {
		slice1 = append(slice1, i)
	}
	for i := 0; i < len(slice1); i++ {
		if slice1[i]%3 != 0 {
			slice2 = append(slice2, slice1[i])
		}
	}
	slice2 = append(slice2, 114514)
	fmt.Println(slice2)
}
