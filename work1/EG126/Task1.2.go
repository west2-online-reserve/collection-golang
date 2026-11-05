//2.洛谷P1046：https://www.luogu.com.cn/problem/P1046

package main

import "fmt"

func main() {

	var arr [10]int
	var a, b, j int
	b = 0

	for i := 0; i < 10; i++ {
		fmt.Scan(&arr[i])
	}

	fmt.Scan(&a)

	for _, j = range arr {
		if j <= (a + 30) {
			b++
		}
	}

	fmt.Println(b)
}
