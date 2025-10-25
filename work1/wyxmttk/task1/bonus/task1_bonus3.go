package main

func findSum(arr []int, target int) [2]int {
	var table = make(map[int]int)
	for i, v := range arr {
		need := target - v
		index, b := table[need]
		if !b {
			table[v] = i
			continue
		} else {
			return [2]int{i, index}
		}
	}
	return [2]int{-1, -1}
}
