package main

import "fmt"

func main() {
	var slice1 []int
	var i int
	for i = 0; i < 50; i++ {
		slice1 = append(slice1, i+1)
	}
	var slice2 []int
	for i = 0; i < 50; i++ {
		if slice1[i]%3 == 0 {
			slice2 = append(slice2, slice1[i])
		}
	}
	slice2 = append(slice2, 114514)
	fmt.Println(slice2)

}
