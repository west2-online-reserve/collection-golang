package main

import "fmt"

//go使用defer＋recover来处理错误
func main() {
	test()
	fmt.Println("执行成功")
}

// Go 的 panic 用于处理不可恢复的错误，recover 用于从 panic 中恢复。
// panic:
// 导致程序崩溃并输出堆栈信息。
// 常用于程序无法继续运行的情况。
// recover:
// 捕获 panic，避免程序崩溃。
func test() {
	defer func() {
		//调用recover内置函数，可以捕获错误
		err := recover()
		//假如没有捕获到错误，放回零值nil
		if err != nil {
			fmt.Println("错误已经捕获")
			fmt.Println("err是", err)
		}
	}() //加括号表示匿名函数被调用
	num1 := 10
	num2 := 0
	result := num1 / num2
	fmt.Println(result)
}
