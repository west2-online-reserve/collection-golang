package main

import "fmt"

func main() {
	var x, y int
	var a [1000]int
	cnt := 0
	fmt.Scanf("%d%d", &x, &y)
	for i := x; i <= y; i++ {
		if (i%4 == 0 && i%100 != 0) || (i%400 == 0) {
			a[cnt] = i
			cnt++
		}
	}
	fmt.Printf("%d\n", cnt)
	for i := 0; i < cnt; i++ {
		fmt.Printf("%d ", a[i])
	}
}

