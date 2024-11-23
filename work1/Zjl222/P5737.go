package main

import "fmt"

func main() {
	var a int
	var b int

	fmt.Scan(&a)
	fmt.Scan(&b)

	s1 := make([]int,1500,1500)
	cnt:=0
	for i := a; i <= b; i++ {
		if (i%400 == 0) || (i%100 != 0 && i%4 == 0) {
			s1[cnt]=i
			cnt++
		}
	}
	fmt.Println(cnt)
	
	for i:=0;i<cnt;i++{
		fmt.Printf("%d ",s1[i])
	}
	
}
