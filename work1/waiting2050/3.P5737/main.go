package main

import "fmt"

func isLeap(y int) bool {
	return (y % 4 == 0 && y % 100 != 0) || y % 400 == 0
}

func main() {
	
	var x, y int
	fmt.Scan(&x, &y)

	var cnt int
	v := []int{}
	for i := x; i <= y; i++ {
		if isLeap(i) {
			cnt++
			v = append(v, i)
		}
	}

	fmt.Println(cnt)
	for _, val := range v {
		fmt.Printf("%d ", val)
	}
}
