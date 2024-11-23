package main

import "fmt"

func main() {

	var s1 = make([]int, 10, 10)
	var x int
	var flag int
	cnt:=0
	for i := 0; i < 10; i++ {
		fmt.Scan(&x)
		s1[i] =x
	}

	fmt.Scan(&flag)
	
	for i := 0; i < 10; i++ {
		if s1[i]<=flag+30 {
			cnt++
		}
	}
	
	fmt.Println(cnt)
}