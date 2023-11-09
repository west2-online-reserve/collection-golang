package main

import "fmt"

func twoSum(nums []int, target int) []int {
	m := make(map[int]int)
	for k, val := range nums {
		_, ok := m[target-val]
		if ok {
			return []int{k, m[target-val]}
		}
		m[val] = k
	}
	return nil
}
func main() {
	var n int
	var target int
	fmt.Scan(&n, &target)
	num := []int{}
	for i := 0; i < n; i++ {
		var a int
		fmt.Scan(&a)
		num = append(num, a)
	}
	fmt.Printf("twoSum(nums, target): %v\n", twoSum(num, target))
}
