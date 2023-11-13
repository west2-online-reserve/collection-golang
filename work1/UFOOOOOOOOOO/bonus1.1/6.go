package main

import (
	"fmt"
	
)
func main(){
	var num1 [9]int = [9]int{1,2,3,4,5,6,7,8,9}
	var num2 [9]int = [9]int{1,2,3,4,5,6,7,8,9}
	for a := 0 ; a <9 ; a++ {
		layout(a,num1,num2)
	}
}
func layout(x int,num1 [9]int,num2 [9]int){
	for b := 0 ; b < 9 ; b++ {
		c := num1[x]*num2[b]
		fmt.Print(x+1,"*",b+1,"=",c)
		fmt.Print(" ")
		if b == 8 {
			fmt.Println()
		}
		if b == 9 {
			break
		}
	 	
	}
}