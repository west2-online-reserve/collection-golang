package main

func main() {
	nums1 := []int{2, 7, 11, 15}
	target1 := 9
	result1 := twoSum(nums1, target1)
	println(result1[0], result1[1])

	nums2 := []int{3, 2, 4}
	target2 := 6
	result2 := twoSum(nums2, target2)
	println(result2[0], result2[1])
}

func twoSum(nums []int, target int) []int {
	numsMap := make(map[int]int)
	for i, num := range nums {
		if v, exist := numsMap[target-num]; exist {
			return []int{v, i}
		}
		numsMap[num] = i
	}
	return nil
}
