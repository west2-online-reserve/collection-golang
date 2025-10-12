package main

import "fmt"

//name...type 表示可以传入任意个type类型的参数
//在函数体内，name当做是一个切片使用
func test(args ...int) {
	fmt.Println("args=", args)
	for index, value := range args {
		fmt.Printf("args[%v]=%v\n", index, value)
	}
}

func test2(num ...*int) {
	for index, value := range num {
		fmt.Printf("%d,%v,%v\n", index, value, *value)
	}
}

func main() {
	test(10, 20, 30, 40, 50)
	var num1, num2, num3 int = 1, 2, 3
	test2(&num1, &num2, &num3)
}
