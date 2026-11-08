package main

import "fmt"

func main() {
	var slice []int
	for i := 1; i <= 50; i++ {
		slice = append(slice, i)
	}
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i]%3 == 0 {
			slice = append(slice[:i], slice[i+1:]...)
		}
	}
	slice = append(slice, 114514)
	fmt.Println(slice)
}
