package main

import "fmt"

func main() {
	var x, y int
	fmt.Scanf("%d %d", &x, &y)
	n := 0
	arr := []int{}
	for i := x; i <= y; i++ {
		if i%4 == 0 && i%100 != 0 || i%400 == 0 {
			n++
			arr = append(arr, i)
		}
		continue
	}
	fmt.Println(n)
	for _, num := range arr {
		fmt.Print(num, " ")
	}
}
