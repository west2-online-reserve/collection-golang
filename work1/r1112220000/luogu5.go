package main

import "fmt"

func main() {
	var M []int
	for i := 1; i <= 50; i++ {
		if i%3 != 0 {
			M = append(M, i)
		}
	}
	M = append(M, 114514)
	fmt.Println(M)
}
