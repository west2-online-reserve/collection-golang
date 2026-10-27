package main

import "fmt"

func main() {
	var new []int
	for i := 1; i <= 50; i++ {
		if i % 3 == 0 {
			continue
		}
		new = append(new, i)
	}
	new = append(new, 114514)
	fmt.Println(new)
}
