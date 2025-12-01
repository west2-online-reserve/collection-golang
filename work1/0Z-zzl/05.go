package main

import "fmt"

func main() {
	var slice []int
	for i := 1; i <= 50; i++ {
		slice = append(slice, i)
	}
	var newslice []int
	for _, v := range slice {
		if v%3 != 0 {
			newslice = append(newslice, v)
		}
	}
	newslice = append(newslice, 114514)
	fmt.Println(newslice)
}
