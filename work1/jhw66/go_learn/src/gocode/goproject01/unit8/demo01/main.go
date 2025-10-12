package main

import "fmt"

func main() {
	var a map[int]string = make(map[int]string, 10) //map可以存放10个键值对
	//or a := make(map[int]string, 10) // map can store 10 key-value pairs
	//只声明map内存是没有分配空间的
	//必须通过make函数进行初始化，才会分配空间

	a[2009] = "hhh"
	a[2008] = "lll"
	a[2007] = "ppp"
	a[2009] = "hhh1" //会替代前面的
	fmt.Println(a)

	b := map[int]string{
		2009: "hhh",
		2008: "lll",
	}
	b[2000] = "kkk"

	fmt.Print(b)
}
