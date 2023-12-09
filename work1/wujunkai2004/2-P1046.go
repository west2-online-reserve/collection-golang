package main

import "fmt"

func main()  {
	var height int = 0
	var apples[16] int
	var ans int =0
	
	for idx :=0; idx<10; idx++{
		fmt.Scan( &apples[idx])
	}

	fmt.Scan( &height)

	for idx := 0; idx<10; idx++ {
		if( apples[idx] <= height + 30 ){
			ans++;
		}
	}

	fmt.Printf("%d", ans)
}