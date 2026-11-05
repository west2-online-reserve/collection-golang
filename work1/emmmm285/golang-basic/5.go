package main

import (
	"fmt"
	"slices"
)

func main() {
	s := make([]int, 50)
	for i := 0; i < 50; i++ {
		s[i] = i + 1
	}
	for i := 0; i < len(s); i++ {
		if s[i]%3 == 0 {
			s = slices.Delete(s, i, i+1)
			i--
		}
	}
	s = append(s, 114514)
	for i := 0; i < len(s); i++ {
		if i != len(s)-1 {
			fmt.Printf("%d ", s[i])
		} else {
			fmt.Printf("%d ", s[i])
		}
	}
}
