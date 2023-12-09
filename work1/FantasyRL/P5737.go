package main

import "fmt"

func main() {
	var x, y int
	fmt.Scanf("%d%d", &x, &y)
	t := y - x
	var yr [1500]int //缓冲区
	var ryr []int
	var count int
	for i := 0; i <= t; i++ {
		yr[i] = x + i
		if (yr[i]%4 == 0 && yr[i]%100 != 0) || yr[i]%400 == 0 { //(yr[i]%4 == 0 && yr[i]%100 != 0) || yr[i]%400 == 0
			ryr = append(ryr, yr[i])
			count++
		}
	}
	fmt.Println(count)
	for i := 0; i < count; i++ {
		fmt.Printf("%d ", ryr[i])
	}
}
