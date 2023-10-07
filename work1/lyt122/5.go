package main

import "fmt"

func task5() []int {
	slice := make([]int, 0)
	for i := 1; i <= 50; i++ {
		slice = append(slice, i)
	}

	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i]%3 == 0 {
			slice = append(slice[:i], slice[i+1:]...)
		}
	}
	slice = append(slice, 114514)
	return slice
}
func main() {
	slice := task5()
	for _, i := range slice {
		fmt.Printf("%d ", i)
	}
}
