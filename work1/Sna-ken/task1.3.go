package main

import "fmt"

func main() {
	var x, y, cnt int
	var sli []int

	fmt.Scanln(&x, &y)

	for ; x <= y; x++ {
		if x%400 == 0 || x%4 == 0 && x%100 != 0 {
			sli = append(sli, x)
			cnt++
		}
	}
	fmt.Println(cnt)

	for i, num := range sli {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(num)
	}

}
