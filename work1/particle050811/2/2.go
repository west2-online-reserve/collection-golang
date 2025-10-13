package main

import "fmt"

func main() {
	a := make([]int, 10)
	for i := 0; i < 10; i++ {
		fmt.Scan(&a[i])
	}
	var len int
	fmt.Scan(&len)

	cnt := 0

	for i := 0; i < 10; i++ {
		if len+30 >= a[i] {
			cnt++
		}
	}
	fmt.Println(cnt)
}
