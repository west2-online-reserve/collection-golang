package main

import "fmt"

func main() {
	var x, y int
	var a [1419]int
	num, index := 0, 0
	fmt.Scanf("%d %d", &x, &y)
	for i := x; i <= y; i++ {
		if i%400 == 0 || i%4 == 0 && i%100 != 0 {
			num += 1
			a[index] = i
			index += 1
		}
	}
	fmt.Printf("%d\n", num)
	for i := 0; i < index; i++ {
		fmt.Printf("%d ", a[i])
	}
}
