package main

import (
	"fmt"
)

func main() {
	list := make([]int, 50)
	for i := 1; i <= 50; i++ {
		list[i-1] = i
	}

	for i := 0; i < len(list)-1; i++ {
		if list[i]%3 == 0 {
			list = append(list[:i], list[i+1:]...)
		}
	}
	list = append(list, 114514)
	fmt.Println(list)
}
