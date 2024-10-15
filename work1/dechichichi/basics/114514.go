package main

import "fmt"

func main() {
	sli := make([]int, 50)
	for i := 0; i < 50; i++ {
		sli[i] = i + 1
	}
	sli2 := make([]int, 0)
	for _, i := range sli {
		if i%3 != 0 {
			sli2 = append(sli2, i)
		}
	}
	sli2 = append(sli2, 114514)
	sli = sli2
	fmt.Println(sli)
}
