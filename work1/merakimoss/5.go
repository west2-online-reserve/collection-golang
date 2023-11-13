package main

import "fmt"

func main() {
	slc := make([]int, 50)
	ans := []int{}
	for x := range slc {
		slc[x] = x + 1
	}
	//	fmt.Println(slc)
	for x := range slc {
		if slc[x]%3 != 0 {
			ans = append(ans, slc[x])
		}
	}
	ans = append(ans, 666) //666
	fmt.Println(ans)
}
