package main

import "fmt"

//批量声明全局变量
var (
	n9  = 500
	n10 = "mary"
)

func main() {
	var num1 int = 10
	fmt.Println(num1)

	var num2 int
	fmt.Println(num2)

	//自行判断变量类型
	var num3 = "tom"
	fmt.Println(num3)

	//使用:=自动推导变量类型
	sex := "男"
	fmt.Println(sex)

	//一次性声明多个变量
	var n1, n2, n3 = 10, "hack", 7.8
	n4, n5, n6 := 100, "jack", 6.7
	fmt.Println(n1, n2, n3)
	fmt.Println(n4, n5, n6)

	fmt.Println(n9, n10)

}
