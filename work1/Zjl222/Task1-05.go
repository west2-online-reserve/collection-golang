package main

import "fmt"
func main(){

	numbers:=make([]int,50)
	for i:=0;i<50;i++{
		numbers[i]=i+1
	}


	answers:=make([]int,0,len(numbers))

	for _,v:=range numbers{
		if v % 3!=0 {
			answers=append(answers,v)
		}
	}
	answers = append(answers,114514)
	
	fmt.Println(answers)
}