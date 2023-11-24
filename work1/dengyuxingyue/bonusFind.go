package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// 使用哈希表来实现查找
func twoSum(nums []int, target int) []int {
	numMap := make(map[int]int)

	for i, num := range nums {
		diff := target - num
		if j, ok := numMap[diff]; ok && j != i {
			return []int{j, i}
		}
		numMap[num] = i
	}

	return []int{}
}

func main() {

	//以下代码均为处理数据输入的代码，实现输入数字，空格区分，换行终止，得到数组的功能

	var numint []int    //我们的目标数组
	var str string = "" //过度字符串

	fmt.Println("请输入你想要的整数数组,换行代表结束")

	reader := bufio.NewReader(os.Stdin)

	nums, _, err := reader.ReadLine()

	if nil != err {
		fmt.Println("reader.ReadLine() error:", err)
	}

	for i := 0; i < len(nums); i++ {

		if int(nums[i]) != 32 {
			str += string(nums[i])
		}
		if int(nums[i]) == 32 || i == len(nums)-1 {

			num, err := strconv.Atoi(str)
			if err == nil {
				fmt.Println("正在加载数据：", num)
			} else {
				fmt.Println("error")
			}
			numint = append(numint, num)
			str = ""
		}
	}
	//输入功能结束

	fmt.Println("加载完毕，请输入你想要的目标值")
	var target int
	fmt.Scan(&target)

	fmt.Println(twoSum(numint, target))

}
