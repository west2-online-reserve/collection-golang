package main

import "fmt"

func main() {
	var a, b, c, d, e, f, g, h, i, j int
	var t int

	fmt.Scanln(&a, &b, &c, &d, &e, &f, &g, &h, &i, &j)

	sli := []int{a, b, c, d, e, f, g, h, i, j}
	var cnt int
	fmt.Scanln(&t)

	for n := 0; n < len(sli); n++ {
		if sli[n] <= t+30 {
			cnt++
		}
	}

	fmt.Println(cnt)
}
