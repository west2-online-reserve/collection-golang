package main

import "fmt"

func isLeapyear(x int) bool {
	if x%400 == 0 {
		return true
	}
	if x%100 == 0 {
		return false
	}
	if x%4 == 0 {
		return true
	}
	return false
}
func main() {
	res := make([]int, 0)
	var l, r int
	fmt.Scan(&l, &r)
	for i := l; i <= r; i++ {
		if isLeapyear(i) {
			res = append(res, i)
		}
	}
	fmt.Println(len(res))
	for _, x := range res {
		fmt.Printf("%d ", x)
	}

}
