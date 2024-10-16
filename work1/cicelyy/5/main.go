package main

import "fmt"

func main() {
	numbers := make([]int, 50)
	for i := range numbers {
		numbers[i] = i + 1
	}
	j := 0
	for i := 0; i < len(numbers); i++ {
		if numbers[i]%3 != 0 {
			numbers[j] = numbers[i]
			j++
		}
	}
	numbers = numbers[:j]
	numbers = append(numbers, 114514)
	fmt.Println(numbers)
}
