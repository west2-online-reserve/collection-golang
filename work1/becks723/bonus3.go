package main

import "fmt"

func main() {
	fmt.Println(sumOfTwoNums([]int{2, 7, 11, 15}, 9))
	fmt.Println(sumOfTwoNums([]int{3, 2, 4}, 6))
}

/* 两数之和 - O(n) */
func sumOfTwoNums(nums []int, target int) []int {
	m := make(map[int]int)
	for i, v := range nums {
		if j, ok := m[target-v]; ok {
			return []int{j, i} // i后被记录因此先输出j
		}
		m[v] = i
	}
	return nil
}
