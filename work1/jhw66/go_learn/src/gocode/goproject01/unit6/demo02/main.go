package main

import "fmt"

func main() {
	//定义二维数组
	var arr [2][3]int16 = [2][3]int16{{1, 4, 7}, {2, 5, 8}}
	fmt.Println(arr)
	fmt.Printf("%p %p %p", &arr, &arr[0], &arr[0][0])

	//赋值：
	arr[0][1] = 7
	arr[0][0] = 82
	arr[1][1] = 25
	fmt.Println(arr)

	//二维数组遍历
	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr[i]); j++ {
			fmt.Print(arr[i][j], "\t")
		}
		fmt.Println()
	}

	fmt.Println("----------------------------")

	for key, value := range arr {
		for k, v := range value {
			fmt.Printf("arr[%v][%v]=%v\t", key, k, v)
		}
		fmt.Println()
	}
}
