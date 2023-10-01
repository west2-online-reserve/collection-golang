package main

import "fmt"

func main() {
	var start, end, num int
	fmt.Scan(&start, &end)
	for i := start; i <= end; i++ {
		if (i%4 == 0 && i%100 != 0) || i%400 == 0 {
			num++
		}
	}
	fmt.Println(num)
	for i := start; i <= end; i++ {
		if (i%4 == 0 && i%100 != 0) || i%400 == 0 {
			fmt.Printf("%d ", i)
		}
	}
}
