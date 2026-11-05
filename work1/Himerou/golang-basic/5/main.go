package main

import "fmt"

func main() {
	s := make([]int, 50)
	for i, _ := range s {
		s[i] = i + 1
	}
	for i, v := range s {
		if v%3 == 0 {
			s = append(s[:i], s[i+1:]...)
		}
	}
	s = append(s, 114514)
	fmt.Println(s)
}
