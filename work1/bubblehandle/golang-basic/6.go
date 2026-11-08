package main

import "fmt"

func main() {
	var arr [9]int
	for i := 0; i < 9; i++ {
		arr[i] = i + 1
	}
	for i := 0; i < 9; i++ {
		var Arr []int
		for j := i; j < 9; j++ {
			Arr = append(Arr, arr[i]*arr[j])
			fmt.Printf("%-4d", Arr[j-i])
		}
		fmt.Println()
	}
}
