package main

import "fmt"

//将匿名函数给一个全局变量
var mul = func(num1 int, num2 int) int {
	return num1 * num2
}

func Calculate(num1 int, num2 int, operation func(int, int) int) int {
	return operation(num1, num2)
}

func multiplier(factor int) func(int) int {
	return func(x int) int {
		return x * factor
	}
}

// 在 Go 语言中，匿名函数是指没有函数名的函数，也常被称为函数字面量。它们可以直接定义在代码中，并像普通值一样被传递和使用。
func main() {
	//定义匿名函数时直接调用，这个方式匿名函数只能使用一次
	result := func(num1 int, num2 int) int {
		return num1 + num2
	}(10, 20)

	//将匿名函数赋给一个变量，这个变量实际就是这个函数类型的变量
	sub := func(num1 int, num2 int) int {
		return num1 - num2
	}
	result2 := sub(30, 10)
	result3 := mul(20, 30)
	//匿名函数作为参数传递
	result4 := Calculate(20, 40, func(num1 int, num2 int) int {
		return num1 + num2
	})
	//匿名函数作为返回值
	result5 := (multiplier(5))(3)

	fmt.Println(result)
	fmt.Println(result2)
	fmt.Println(result3)
	fmt.Println(result4)
	fmt.Println(result5)
}
