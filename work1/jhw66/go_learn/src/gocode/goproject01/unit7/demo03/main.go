package main

import "fmt"

func main() {
	//切片定义以后不能直接使用，需要将其引用到一个数组，或者make一个空间供切片使用
	var slice []int
	fmt.Println(slice)

	var intarr [6]int = [6]int{1, 4, 7, 2, 5, 8}
	slice = intarr[:] //等于是slice:=intarr[0:len(arr)]
	fmt.Println(slice)

	//切片再次切片
	slice2 := slice[1:]
	fmt.Println(slice2)
	slice2[0] = 66
	fmt.Println(intarr)
	fmt.Println(slice)

	//切片后追加元素
	slice = append(slice, 88, 50)
	fmt.Println(slice)
	// 	如果 len < cap：直接放到原底层数组后面；
	// 如果 len == cap：Go 会自动创建一个更大的底层数组，然后把旧数据复制过去。

	slice3 := []int{99, 44}
	slice3 = append(slice, slice3...) //...表示追加的是一个切片
	fmt.Println(slice3)
	slice = make([]int, 3, 10)
	slice = append(slice, slice3...)
	fmt.Println(slice)

	//切片的拷贝
	var slice4 []int = make([]int, 3)
	copy(slice4, slice3) //将slice3中对应数组中的元素内容复制到slice4中,受slice4大小限制
	fmt.Println(slice4)
}
