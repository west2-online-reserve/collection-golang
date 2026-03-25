package main

import "fmt"

func main() {
	var Apples [10]int
	for i := 0; i < 10; i++ {
		fmt.Scan(&Apples[i])
	}

	var Reach int
	fmt.Scan(&Reach)

	MaxReach := Reach + 30
	Counts := 0

	for _, apple := range Apples {
		if apple <= MaxReach {
			Counts++
		}
	}

	fmt.Println(Counts)
}
