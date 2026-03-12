//3.洛谷P5737：https://www.luogu.com.cn/problem/P5737

package main

import "fmt"

func Isleap(year int) bool {
	if year%400 == 0 || (year%4 == 0 && year%100 != 0) {
		return true
	}
	return false
}

func main() {
	var x, y, num, year int
	num = 0
	fmt.Scan(&x, &y)
	arr := []int{}
	for i := x; i <= y; i++ {
		if Isleap(i) {
			arr = append(arr, i)
			num++
		}
	}
	fmt.Println(num)
	for _, year = range arr {
		fmt.Printf("%d ", year)
	}
}
