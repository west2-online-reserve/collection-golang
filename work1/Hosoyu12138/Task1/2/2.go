package main

import "fmt"

func main() {
	a := 30
	var height int
	var arr [10]int
	var sum int
	fmt.Scanf("%d%d%d%d%d%d%d%d%d%d\n", &arr[0], &arr[1], &arr[2], &arr[3], &arr[4], &arr[5], &arr[6], &arr[7], &arr[8], &arr[9])
	fmt.Scanf("%d", &height)
	for i := 0; i <= 9; i++ {
		if arr[i] <= (height + a) {
			sum += 1
		}
	}
	fmt.Println(sum)
}
