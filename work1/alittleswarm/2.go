package main

import "fmt"

func main() {
	var a [10]int
	var Len, cnt = 0, 0

	for i := 0; i < 10; i++ {
		fmt.Scanf("%d", &a[i])
	}
	fmt.Scanf("%d", &Len)
	for _, j := range a {
		if Len+30 > j {
			cnt++
		}

	}
	print(cnt)
}
