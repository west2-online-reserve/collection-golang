package main

import "fmt"

func main() {
	map1 := make(map[int]int)
	var arr = [4]int{2, 7, 11, 15} //初始数组
	var slice []int                //用于存储下标的切片
	var target int = 9             //目标和
	var tempt int                  //用于存储 目标和 和数组中的值的差值
	for i := 0; i < len(arr); i++ {
		map1[arr[i]] = 0
	}

	for i := 0; i < len(arr); i++ {
		tempt = target - arr[i]
		if _, isOk := map1[tempt]; isOk {
			slice = append(slice, i)
			break

		}

	}
	for i := 0; i <= (len(arr) - 1); i++ {
		if arr[i] == tempt {
			slice = append(slice, i)
		}

	}
	fmt.Println(slice)
}

//易于想到的方法是利用嵌套for循环遍历整个数组，该方法是通过查询搜索引擎上的方法后，根据其使用map进行寻找数值的方法重新写过的
