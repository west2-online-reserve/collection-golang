package	main

import "fmt"

func main(){
	
	a := [11]int{}
	for i := 1; i <= 10; i++{
		fmt.Scan(&a[i])
	}

	var base, cnt int
	fmt.Scan(&base)
	for i := 1; i <= 10; i++{
		if base + 30 >= a[i]{
			cnt++
		}
	}

	fmt.Println(cnt)
}