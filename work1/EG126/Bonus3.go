//3.给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出 和为目标值 target 的那 两个 整数，并返回它们的数组下标。

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func twoSum(nums []int, target int) (int, int) {
	numMap := make(map[int]int)
	for i, val := range nums {
		complement := target - val
		if j, exists := numMap[complement]; exists {
			return j, i
		}
		numMap[val] = i
	}
	return -1, -1
}

func main() {
	scanner := bufio.NewScanner(os.Stdin) // 用于读取输入，避免缓冲区问题

	fmt.Print("输入数组: ")
	scanner.Scan()
	parts := strings.Fields(scanner.Text()) // 按空格分割输入字符串
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p) // 字符串转整数
	}

	fmt.Print("输入目标值: ")
	scanner.Scan()
	target, _ := strconv.Atoi(scanner.Text()) // 字符串转整数

	i, j := twoSum(nums, target)
	fmt.Printf("[%d,%d]\n", i, j)
}
