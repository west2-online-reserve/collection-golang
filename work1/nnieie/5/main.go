package main

import "fmt"

func main() {
    s := make([]int, 0, 50)
    for i := 1; i <= 50; i++ {
        s = append(s, i)
    }
	for i := 0; i < len(s); i++ {
		if s[i]%3 == 0 {
			s = append(s[:i], s[i+1:]...)
			i--
		}
	}
	s = append(s, 114514)
    fmt.Println(s)
}