package main

import "fmt"

func runnian(i int) bool {
	var ret = false
	if i%4 == 0 && i%100 != 0 || i%400 == 0 {
		ret = true
	}
	return ret
}
func main() {
	var x, y, i int
	var cnt = 0
	var slice []int
	fmt.Scanf("%d %d", &x, &y)
	for i = x; i <= y; i++ {
		if runnian(i) {
			cnt++
			slice = append(slice, i)
		}
	}
	fmt.Println(cnt)
	fmt.Println(&slice)

}
