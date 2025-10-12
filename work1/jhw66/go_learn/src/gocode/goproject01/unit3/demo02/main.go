package main

import "fmt"

func main() {
	var age int
	fmt.Print("请输入年龄：")
	fmt.Scanln(&age)
	fmt.Println("你输入的年龄是：", age)

	var name string
	fmt.Print("请输入姓名：")
	fmt.Scanf("%s", &name)
	fmt.Println("你输入的姓名是：", name)
}
