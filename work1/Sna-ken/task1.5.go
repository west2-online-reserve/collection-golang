package main

import "fmt"

func main() {
	sli := make([]int, 0)
	for i := 1; i <= 50; i++ {
		sli = append(sli, i)
	}

	for i := len(sli) - 1; i >= 0; i-- {
		if sli[i]%3 == 0 {
			sli = append(sli[:i], sli[i+1:]...)
		}
	}

	sli = append(sli, 114514)

	fmt.Println(sli)
}
