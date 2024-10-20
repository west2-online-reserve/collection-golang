package main

import "fmt"

func main() {
	var x int
	var y int
	fmt.Scanf("%d", &x)
	fmt.Scanf("%d", &y) // 读取y的值
	cnt := 0
	var list []int
	for i := x; i <= y; i++ { // 确保循环包括y
		if isRun(i) {
			list = append(list, i) // 使用append来添加元素
			cnt++
		}
	}
	fmt.Println(cnt)
	for _, value := range list {
		fmt.Printf("%d ", value) // 使用value而不是list
	}
	fmt.Println() // 打印换行
}

func isRun(year int) bool {
	if year%400 == 0 || (year%100 != 0 && year%4 == 0) {
		return true
	}
	return false
}
