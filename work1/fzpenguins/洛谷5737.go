package main

import "fmt"

func main() {
	var x, y int
	fmt.Scan(&x, &y)
	var runNian []int
	cnt := 0
	for i := x; i <= y; i++ {
		if (i%4 == 0 && i%100 != 0) || (i%400 == 0) {
			runNian = append(runNian, i)
			cnt++
		}
	}
	fmt.Println(cnt)
	for i := 0; i < cnt; i++ {
		fmt.Print(runNian[i], " ")
	}
}
