//自己只想到了o(n^2)的解法(对手指),
// 另外在输入数据部分有参考豆包给的方案

package main

import "fmt"

func main() {
	var num int
	nums := make([]int, 0)
	sca := make([]int, 0)
	for {
		_, err := fmt.Scan(&num)
		if err != nil {
			break
		}
		sca = append(sca, num)
	}
	nums = sca[:len(sca)-1]
	target := sca[len(sca)-1:]
	r := fun1(nums, target[0])
	fmt.Println(r)

}

func fun1(n []int, t int) []int {
	for i := 0; i < len(n); i++ {
		for j := 0; j < len(n); j++ {
			sum := n[i] + n[j]
			if sum == t && i != j {
				return []int{i, j}
			}
		}
	}
	return []int{}
}
