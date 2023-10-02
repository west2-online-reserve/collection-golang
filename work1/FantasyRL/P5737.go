package main

import "fmt"

func main() {
	var x, y int
	fmt.Scanf("%d%d", &x, &y)
	var ryr []int
	var count int
	for i := x; i <= y; i++ {
		if (i%4 == 0 && i%100 != 0) || i%400 == 0 { //(yr[i]%4 == 0 && yr[i]%100 != 0) || yr[i]%400 == 0
			count++
			ryr = append(ryr, i)

		}
	}
	fmt.Println(count)
	for i := 0; i < count; i++ {
		fmt.Printf("%d ", ryr[i])
	}
}
