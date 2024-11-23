package main

import "fmt"

func main(){
 	nums:=make(map[int]int)
	var n,target int
	fmt.Scan(&n)
	for i:=0;i<n;i++ {
		var t int
		fmt.Scan(&t)
		nums[t] = i
	}
	fmt.Scan(&target)
	for num,x := range nums {      //num是数组元素，x是数组下标
		y,ok := nums[target-num]       //y是满足条件的另一个数组下标
		if(ok&&y!=x){                  //排除数组中同一个元素在答案里重复出现的情况
			fmt.Printf("[%d,%d]",x,y)
			break
		}
	}
}