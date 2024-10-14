package main

import (
	"fmt"
	"sort"
)

func main() {
	nums := []int{2, 7, 11, 15}
	map1 := make(map[int]int)
	target := 9

	for i, v := range nums {
		map1[v] = i
	}

	sort.Ints(nums)

	for i, j := 0, len(nums)-1; i < j && i < len(nums); i++ {
		for i < j && nums[i]+nums[j] > target {
			j--
		}

		if i != j && nums[i]+nums[j] == target {
			fmt.Println([]int{map1[nums[i]], map1[nums[j]]})
			break
		}
	}
}
