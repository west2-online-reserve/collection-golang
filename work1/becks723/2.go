package main

import "fmt"

func main() {
	var a [10]int
	var ans, h int
	for i := 0; i < len(a); i++ {
		fmt.Scanf("%d", &a[i])
	}
	fmt.Scanf("%d", &h)
	for _, v := range a {
		if v <= h+30 {
			ans++
		}
	}
	fmt.Println(ans)
}
