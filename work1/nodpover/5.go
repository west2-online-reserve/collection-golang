package main

import "fmt"

func main() {
	s1 := make([]int, 50)
	for i := 1; i <= 50; i++ {
		s1[i-1] = i
	}

	for i := 0; i <= len(s1)-1; i++ {
		if s1[i]%3 == 0 {
			s1 = append(s1[:i], s1[i+1:]...)
		}
	}

	s1 = append(s1, 114514)
	fmt.Println(s1)
}
