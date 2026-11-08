package main

import "fmt"

func main() {
	var a, b, ans int
	fmt.Scan(&a, &b)
	var arr []int
	for i := a; i <= b; i++ {
		if i%400 == 0 || (i%4 == 0 && i%100 != 0) {
			arr = append(arr, i)
			ans++
		}
	}
	fmt.Println(ans)
	for _, year := range arr {
		fmt.Print(year, " ")
	}
}
