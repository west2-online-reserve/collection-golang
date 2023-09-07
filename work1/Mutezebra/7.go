package main

import "fmt"

func main() {
	nums := []int{2, 7, 4, 9}
	target := 6

	result, ok := solve(nums, target)
	if ok {
		fmt.Println(result)
	} else {
		fmt.Println("Not exist")
	}
}

func solve(nums []int, target int) ([]int, bool) {
	m := make(map[int]int)

	for k, v := range nums {
		if value, ok := m[target-v]; ok {
			return []int{value, k}, true
		} else {
			m[v] = k
		}
	}
	return nil, false
}
