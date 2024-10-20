package main

import "fmt"

func main() {
	var a, b, i int
	var num [401]int
	res := 0
	fmt.Scan(&a, &b)
	for i = a; i <= a+8 && i <= b; i++ {
		if run(i) {
			res++
			num[res] = i
			break
		}
	}
	i += 4
	for ; i <= b; i += 4 {
		if run(i) {
			res++
			num[res] = i
		}
	}
	fmt.Println(res)
	for j := 1; j <= res; j++ {
		fmt.Print(num[j], " ")
	}
}
func run(i int) bool {
	if i%100 == 0 && i%400 == 0 {
		return true
	}
	if i%100 != 0 && i%4 == 0 {
		return true
	}
	return false
}
