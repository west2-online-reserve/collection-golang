package main

import "fmt"

func main(){
	var x,y int
	sum:=0
	fmt.Scan(&x,&y)
	var year [1420]int
	for i:=x;i<=y;i++ {
		if i%400==0 {
			year[sum]=i
			sum++
		}else if i%100!=0&&i%4==0 {
			year[sum]=i
			sum++
		}
	}
	fmt.Printf("%d\n",sum)
	for i:=0;i<sum;i++ {
		fmt.Printf("%d ",year[i])
	}

}