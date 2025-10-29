// 给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出 和为目标值 target 的那 两个 整数，并返回它们的数组下标。
// 你可以假设每种输入只会对应一个答案。但是，数组中同一个元素在答案里不能重复出现。
// 你可以按任意顺序返回答案。
// 示例 1：
// 输入：nums = [2,7,11,15], target = 9 输出：[0,1] 解释：因为 nums[0] + nums[1] == 9 ，返回 [0, 1]
// 示例2
// 输入：nums = [3,2,4], target = 6 输出：[1,2]
// 是否有复杂度O(n)的算法？

// // 想到的方法：o(nlogn)
// package main

// import (
// 	"fmt"
// 	"math/rand"
// 	"time"
// )

// type num struct {
// 	Index  int
// 	number int
// }

// func quick_sort(arr []num, low, high int) {
// 	if low >= high {
// 		return
// 	}

// 	pivotIndex := rand.Intn(high-low+1) + low
// 	arr[low], arr[pivotIndex] = arr[pivotIndex], arr[low]
// 	pivot := arr[low].number

// 	lt, i, end := low, low+1, high
// 	for i <= end {
// 		if arr[i].number < pivot {
// 			arr[lt], arr[i] = arr[i], arr[lt]
// 			lt++
// 			i++
// 		} else if arr[i].number > pivot {
// 			arr[i], arr[end] = arr[end], arr[i]
// 			end--
// 		} else {
// 			i++
// 		}
// 	}

// 	quick_sort(arr, low, lt-1)
// 	quick_sort(arr, end+1, high)
// }
// func main() {
// 	rand.Seed(time.Now().UnixNano())
// 	var n int
// 	fmt.Println("请输入给定数组大小:")
// 	fmt.Scan(&n)
// 	nums := make([]num, n)
// 	fmt.Println("请输入给定数组的值:")
// 	for i := 0; i < n; i++ {
// 		fmt.Scan((&nums[i].number))
// 		nums[i].Index = i
// 	}
// 	var target int
// 	fmt.Println("请输入目标值:")
// 	fmt.Scan(&target)

// 	quick_sort(nums, 0, n-1)

//		flag := false
//		for i, j := 0, n-1; i < j; {
//			if nums[i].number+nums[j].number > target {
//				j--
//			} else if nums[i].number+nums[j].number < target {
//				i++
//			} else {
//				flag = true
//				fmt.Println("[", nums[i].Index, ",", nums[j].Index, "]")
//				i++
//			}
//		}
//		if !flag {
//			fmt.Println("没找到对应下标")
//		}
//	}
//

// 看了ai后：o(n)
package main

import (
	"fmt"
)

func TwoSum(nums []int, target int, n int) []int {
	result := make([]int, 0, n)
	m := make(map[int]int, n)
	for index01, value := range nums {
		if index02, ok := m[target-value]; ok {
			result = append(result, index02, index01)
		}
		m[value] = index01
	}
	return result
}
func main() {
	var n int
	fmt.Println("请输入给定数组大小:")
	fmt.Scan(&n)
	nums := make([]int, n)
	fmt.Println("请输入给定数组的值:")
	for i := 0; i < n; i++ {
		fmt.Scan((&nums[i]))
	}
	var target int
	fmt.Println("请输入目标值:")
	fmt.Scan(&target)

	result := TwoSum(nums, target, n)
	if len(result) == 0 {
		fmt.Println("没找到对应下标")
	} else {
		fmt.Printf("结果下标: %v\n", result)
	}
}
