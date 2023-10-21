package main

import "fmt"

func main() {

	var left, right int
	var years []int

	fmt.Scan(&left, &right)

	for i := left; i <= right; i++ {
		if (i%4 == 0 && i%100 != 0) || (i%400) == 0 {
			years = append(years, i)
		}
	}

	fmt.Println(len(years))

	for _, i := range years {
		fmt.Print(i)
		fmt.Print(" ")
	}

}
