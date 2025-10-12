package main

import (
	"errors"
	"fmt"
)

// 调用 panic() 后，Go 会：
// 立即停止当前函数的执行；
// 依次执行所有已注册的 defer；
// 如果没有被 recover() 捕获，程序会崩溃退出，并打印调用栈（stack trace）。
func main() {
	defer func() {
		fmt.Println("结束")
	}()
	err := test()
	if err != nil {
		fmt.Println("自定义错误：", err)
		//出现错误后想让后续代码不再进行，利用builtin包下的内置函数panic
		panic(err)
	}
	fmt.Println("执行成功")
}

// 错误处理的惯用模式
// // 标准模式：返回结果和错误
//
//	func 函数名(参数) (返回类型error) {
//	    // 业务逻辑
//	    if 发生错误 {
//	        return zeroValue, errors.New("错误描述")
//	    }
//	    return 结果, nil  // 成功时返回nil错误
//	}
func test() (err error) {
	num1 := 10
	num2 := 0
	if num2 == 0 {
		//抛出自定义错误
		return errors.New("除数不能为0")
	} else {
		result := num1 / num2
		fmt.Println(result)
		//如果没有错误，返回零值
		return nil
	}
}
