package main

import "fmt"

// 方法一: 暴力解法
// O(n²) O(1)
func twoSum1(nums []int, target int) []int {
    size := len(nums)
    for i := 0; i < size; i++ {
        for j := i+1; j < size; j++ {
            if nums[i] + nums[j] == target {
				return []int{i, j}
            }
        }
    }
    return nil
}

// 方法二: 双向指针(适用于有顺序的数组, 此处假设为从小到大)
// O(n) O(1)
func twoSum2(nums []int, target int) []int {
	left := 0
	right := len(nums)-1
	for left < right {
		if nums[left] + nums[right] < target {
			left++
		} else if nums[left] + nums[right] > target {
			right--
		} else {
			return []int{left, right}
		}
	}
	return nil
}

// 方法三: 哈希表法
// O(n) O(n)
func twoSum3(nums []int, target int) []int {
	size := len(nums)
	hashMap := make(map[int]int, size)
	// 用键来存储num, 值来存下标
	for i := 0; i < size; i++ {
		key, ok := hashMap[target - nums[i]]
		if ok == true {
			return []int{key, i}
		}
		hashMap[nums[i]] = i
	}
	return nil
}


func main() {
	sum1 := twoSum1([]int{1, 3, 4, 7, 9}, 11)
	fmt.Println(sum1)
	sum2 := twoSum2([]int{1, 3, 4, 7, 9}, 11)
	fmt.Println(sum2)
	sum3 := twoSum3([]int{1, 3, 4, 7, 9}, 11)
	fmt.Println(sum3)
}
