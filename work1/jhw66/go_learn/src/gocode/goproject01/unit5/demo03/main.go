package main

import (
	"fmt"
)

func test(num int) int {
	fmt.Println(num)
	return 0
}

// 给函数取别名
type myFunc func(int) int

// 函数也可以作为形参
func testuse(num1 int, num2 float32, testFunc myFunc) {
	testFunc(num1)
	testFunc((int)(num2))
}

// 可以给返回值命名
func test02(num1 int, num2 int) (result01 int, result02 int) {
	result01 = num1 + num2
	result02 = num1 - num2
	return
}

func main() {
	//函数也是一种数据类型，可以赋给一个变量
	a := test
	fmt.Printf("a的类型是:%T,test函数的类型是:%T\n", a, test)
	a(10)
	test(11)
	testuse(10, 3.19, a)

	//自定义数据类型
	type myInt int
	var num1 myInt = 30
	//a(num1)
	//无法传入，因为虽然是别名，但是在go中编译识别的时候还是认为myInt和int不是一种数据类型
	a(int(num1))

	b, c := test02(10, 20)
	fmt.Printf("%d,%d", b, c)
}
