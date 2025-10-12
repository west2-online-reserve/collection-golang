package main

import "fmt"

func main() {
	//定义切片:make函数的三个参数，1切片类型，2切片长度，3切片容量
	slice := make([]int, 4, 20) //在底层创建一个数组，对外不可见，不能直接操作整个数组，只能使用slice间接操作
	fmt.Println(slice)
	fmt.Println(len(slice))
	fmt.Println(cap(slice))
	slice[0] = 66
	slice[3] = 88
	fmt.Println(slice)

	slice2 := []int{1, 4, 7}
	fmt.Println(slice2)
	fmt.Println(len(slice2))
	fmt.Println(cap(slice2))

	//遍历切片
	for i := 0; i < len(slice); i++ {
		fmt.Printf("slice[%v]=%v\t", i, slice[i])
	}

	for key, value := range slice {
		fmt.Printf("slice[%v]=%v\t", key, value)
	}
}
