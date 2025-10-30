package main

import "fmt"

func getMultiplesOfFour(x, y int) {
	var n int = 0
	var slice []int
	for i := x; i <= y; i++ {
		if (i%4 == 0 && i%100 != 0) || i%400 == 0 {
			slice = append(slice, i)
			n = n + 1
		}
	}
	fmt.Println(n)
	fmt.Println(slice)
}
func main() {
	var x, y int
	fmt.Scanln(&x, &y)
	getMultiplesOfFour(x, y)
}
