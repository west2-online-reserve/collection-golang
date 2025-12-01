/*查了些资料找到了o(n)的算法，一些不懂的地方去问了下AI
在函数部分被各种变量搞得头大，也认识到了注释的重要性TT，
走开一会要看半天来理解自己写了些什么*/

package main

import "fmt"

func main() {
	var num int
	nums := make([]int, 0)
	sca := make([]int, 0)
	//输入部分，先输入数组数据再输入target
	for {
		_, err := fmt.Scan(&num)
		if err != nil {
			break
		}
		sca = append(sca, num)
	}

	nums = sca[:len(sca)-1]
	target := sca[len(sca)-1:]
	r := func1(nums, target[0])
	fmt.Println(r)
}

func func1(nums []int, target int) []int {
	m := make(map[int]int) //将下标作为值，元素作为键

	for i, n := range nums {
		rsl := target - n //寻找补数
		v, exist := m[rsl]
		if exist {
			return []int{v, i}
		}
		m[n] = i
	}

	return []int{}
}
