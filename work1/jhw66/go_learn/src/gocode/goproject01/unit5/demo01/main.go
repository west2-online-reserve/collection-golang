package main

import "fmt"

func cal(num1 int, num2 int) int {
	return num1 + num2
}

//go语言中函数不支持函数重载
//如果想实现类似效果，可以使用不同的函数名
// func cal(num1 int) int {
// 	return num1
// }

func cal2(num1 int, num2 int) (int, int) {
	add, sub := num1+num2, num1-num2
	return add, sub
}

func main() {
	fmt.Println(cal(10, 20))
	fmt.Println(cal(30, 40))
	fmt.Println(cal(50, 60))
	sum := cal(70, 80)
	fmt.Println("sum=", sum)
	add, sub := cal2(10, 20)
	fmt.Println("add=", add, "sub=", sub)
	add2, _ := cal2(30, 40) //忽略返回值
	fmt.Println("add2=", add2)
}
