package main

import "fmt"

func isrun(n int) bool {
	if n%4 == 0 && n%100 != 0 || n%400 == 0 {
		return true
	}
	return false
}

func main() {
	var x, y int
	fmt.Scan(&x, &y)
	count := 0
	slice := make([]int, y-x+1)
	for i := x; i <= y; i++ {
		if isrun(i) {
			slice[count] = i
			count++
		}
	}
	fmt.Println(count)
	for i := 0; i < count; i++ {
		fmt.Printf("%d ", slice[i])
	}
}
