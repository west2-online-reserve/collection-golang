package main

import "fmt"

func main() {
	numbers := make([]int, 50)
	//使得切片中元素为1-50
	for i := range numbers {
		numbers[i] = i + 1
	}
	var newnumbers []int
	for _, i := range numbers {
		if i%3 != 0 {
			newnumbers = append(newnumbers, i)
		}
	}
	newnumbers = append(newnumbers, 114514)
	fmt.Println(newnumbers)
}
