package main

import "fmt"

func main() {
	var slice []int
	for i := 1; i <= 50; i++ {
		slice = append(slice, i)

	}
	for i := 0; i <= (len(slice) - 1); i++ {
		tempt := slice[i]
		if tempt%3 == 0 {
			slice = append(slice[:i], slice[i+1:]...)
			i--

		}
	}
	slice = append(slice, 114514)
	fmt.Println(slice)

}
