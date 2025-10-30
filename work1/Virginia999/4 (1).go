// 输入 x,y，输出 [x,y] 区间中闰年个数，并在下一行输出所有闰年年份数字，使用空格隔开。
// 如何创建一个切片，使其为x，y间所有整数呢
package main

import "fmt"

func main() {
	var x, y int
	fmt.Scanf("%d %d", &x, &y)
	start, end := x, y
	count := 0
	//创建一个空切片
	a := []int{}
	//i的起步是start
	for i := start; i <= end; i++ {
		if (i%400 == 0) || (i%4 == 0 && i%100 != 0) {
			count++
			//往切片中加入动态的元素i
			a = append(a, i)
		}
	}
	fmt.Println(count)
	fmt.Println(a)
}
