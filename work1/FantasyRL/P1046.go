package main

import "fmt"

func main() {
	var apple [10]int
	for i := 0; i < 10; i++ {
		fmt.Scanf("%d", &apple[i])
	}
	var pick, p int
	fmt.Scan(&pick)
	var pick2 int = pick + 30
	for i := 0; i < 10; i++ {
		if apple[i] <= pick2 {
			p++
		} else {

		}
	}
	fmt.Println(p)
}
