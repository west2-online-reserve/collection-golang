package main

import "fmt"

func main() {
	var arr = make([]int, 10)
	var high int
	var n int = 0
	for i := 0; i < 10; i++ {
		fmt.Scan(&arr[i])
	}
	fmt.Scanf("%d", &high)

	for i := 0; i < 10; i++ {
		if arr[i] <= high+30 || (arr[i] >= 100 && arr[i] <= 200) || (high >= 100 && high <= 120) {
			n++
		}
	}
	if n == 0 {
		fmt.Println("error")
	} else {
		fmt.Printf("%d", n)
	}
}
