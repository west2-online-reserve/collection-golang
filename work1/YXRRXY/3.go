package main

import "fmt"

var ans int
var years []int

func ansfunc(a, b int) {
	for i:=a; i<=b; i++ {
		if i % 4 == 0 && i % 100 != 0 || i % 400 == 0 {
			years = append(years, i)
			ans++
		}
	}
}

func main() {
	var x, y int
	fmt.Scan(&x, &y) 
	ansfunc(x, y)
	fmt.Println(ans)
	for _, v := range years {
		fmt.Printf("%d ", v)
	}
}