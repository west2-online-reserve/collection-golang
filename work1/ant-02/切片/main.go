package main

import "fmt"

func main() {
	arr := make([]int, 0)
	for i := 1; i <= 50; i++ {
		if i%3 == 0 {
			continue
		}
		arr = append(arr, i)
	}
	arr = append(arr, 114514)
	for _, x := range arr {
		fmt.Print(x, " ")
	}
}