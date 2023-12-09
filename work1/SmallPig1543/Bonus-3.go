package main

func Function(nums []int, target int) []int {
	hash := map[int]int{}
	for i := 0; i < len(nums); i++ {
		if value, ok := hash[target-nums[i]]; ok {
			return []int{value, i}
		}
		hash[nums[i]] = i
	}
	return []int{}
}
