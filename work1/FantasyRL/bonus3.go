package main

import (
	"fmt"
	"sort"
)

func compare(arr []int, target int) (a, b int, err bool) {
	for i, v := range arr {
		if v < target {
			i2 := sort.Search(len(arr), func(i int) bool { return arr[i] >= target-v })
			if i2 == len(arr) || i2 == i {
				continue
			} else {
				return i, i2, true
				break
			}
		}
	}
	return 0, 0, false
}
func main() {
	var arrlen, target int
	var arr []int
	fmt.Scan(&arrlen)
	for i := 0; i < arrlen; i++ {
		var t int
		fmt.Scan(&t)
		arr = append(arr, t)
	}
	fmt.Scan(&target)
	a, b, c := compare(arr, target)
	if c == false {
		fmt.Println("不中嘞哥")
	} else {
		fmt.Println("\n", a, b)
	}

}
