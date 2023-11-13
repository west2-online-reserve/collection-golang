package main

import "fmt"

func main() {
	numbers := make([]int, 0)
	for i := 1; i <= 50; i++ {
		numbers = append(numbers, i)
	}
	nextNumbers := make([]int, 0)
	for _, num := range numbers {
		if num%3 != 0 {
			nextNumbers = append(nextNumbers, num)
		}
	}
	nextNumbers = append(nextNumbers, 114514)
	fmt.Println(nextNumbers)
}
