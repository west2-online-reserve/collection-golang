package main

import "fmt"

func main() {
	var x, y int
	rn := []int{}
	fmt.Scan(&x, &y)
	for i := x; i <= y; i++ {
		if (i%4 == 0 && i%100 != 0) || (i%400 == 0) {
			rn = append(rn, i)
		}
	}
	fmt.Println(len(rn))
	for _, v := range rn {
		fmt.Printf("%d ", v)
	}
}
