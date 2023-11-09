package main

import "fmt"

func main() {
	var s, e, sum int
	fmt.Scan(&s, &e) //输入开始和结束年份
	a := make([]int, 0, 0)
	for i := s; i <= e; i++ {
		if i%400 == 0 {
			a = append(a, i)
			sum++
		} else if i%4 == 0 && i%100 != 0 {
			a = append(a, i)
			sum++
		}
	}

	fmt.Println(sum)
	for i := 0; i < len(a); i++ {
		fmt.Printf("%d ", a[i])
	}
}
