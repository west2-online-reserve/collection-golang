package main

import "fmt"

func main() {
	var S_y, E_y int
	var list []int
	fmt.Scan(&S_y, &E_y)
	for i := S_y; i <= E_y; i++ {
		if i%400 == 0 || (i%4 == 0 && i%100 != 0) {
			list = append(list, i)
		}
	}
	fmt.Println(len(list))
	for i := 0; i < len(list); i++ {
		fmt.Print(list[i], " ")
	}
}
