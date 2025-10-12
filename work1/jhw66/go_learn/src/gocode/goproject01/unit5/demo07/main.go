package main

import (
	"fmt"
)

func add(num1 int, num2 int) int {

	//在Golang中，程序遇到defer关键字，不会立即执行defer后的语句，
	//是将defer后的语句压入一个栈中，继续执行函数后面的语句
	//在函数执行完毕后，从栈中取出语句开始执行
	defer fmt.Println("num1=", num1)
	defer fmt.Println("num2=", num2)
	//不会改变压入栈的语句，相关的值会拷贝进入栈中，不会随函数后面的变化而变化
	num1 += 90
	num2 += 50
	var sum int = num1 + num2
	fmt.Println("sum=", sum)
	return sum
}

//一般defer中放的是关闭某个资源的内容

func main() {
	fmt.Println(add(30, 60))

}
