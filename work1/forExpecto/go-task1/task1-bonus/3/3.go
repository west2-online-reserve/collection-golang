package main

import "fmt"

func twoSum(nums []int, target int) []int {
	two := map[int]int{}
	for i, x := range nums {
		v := target - x
		if _, ok := two[v]; ok {
			return []int{two[v], i}
		}
		two[x] = i
	}
	return []int{}
}

func main() {
	nums := []int{2, 7, 11, 15}
	target := 9
	fmt.Println(twoSum(nums, target))
	nums2 := []int{3, 2, 4}
	target = 6
	fmt.Println(twoSum(nums2, target))
}
