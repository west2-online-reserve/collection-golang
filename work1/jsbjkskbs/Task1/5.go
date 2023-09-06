package main

import (
	"fmt"
)

func main() {
	var slice []int
	for i := 1; i < 51; i++ {
		slice = append(slice, i)
	}
	p := 0
	for i := 0; i < len(slice); i++ {
		if slice[i]%3 != 0 {
			slice[p] = slice[i]
			p++
		}
	}
	slice = slice[:p]
	slice = append(slice, 114514)
	fmt.Print("[")
	for i := 0; i < len(slice)-1; i++ {
		fmt.Printf("%v ", slice[i])
	}
	fmt.Printf("%v]", slice[len(slice)-1])
}
