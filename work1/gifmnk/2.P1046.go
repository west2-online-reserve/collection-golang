package main

import "fmt"

func main(){
	var h [10] int
	var H int
	sum:=0
	for i:=0;i<10;i++ {
		fmt.Scan(&h[i])
	}
	fmt.Scan(&H)
	H+=30
	for i:=0;i<10;i++ {
		if h[i]<=H {
			sum++
		}
	}
	fmt.Printf("%d",sum)
}