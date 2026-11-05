package main

import "fmt"

func main() {
	var len int
	fmt.Println("请输入数组长度：")
	fmt.Scanln(&len)
	nums := make([]int, len)
	fmt.Println("请输入数组元素：")
	for i := 0; i < len; i++ {
		fmt.Scanln(&nums[i])
	}
	var target int
	fmt.Println("请输入目标值：")
	fmt.Scanln(&target)

	HashMap := make(map[int]int)
	for i, num := range nums {
		if j, ok := HashMap[target-num]; ok {
			fmt.Printf("[%d,%d]", i, j)
			return
		}
		HashMap[num] = i
	}
}


