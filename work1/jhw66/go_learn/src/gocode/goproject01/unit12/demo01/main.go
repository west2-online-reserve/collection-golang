package main

import (
	"fmt"
	"reflect"
)

func testReflect(i interface{}) {
	//调用TypeOf函数，返回reflect，Type类型数据
	reType := reflect.TypeOf(i)

	//调用ValueOf函数，返回reflect，Value类型数据
	reValue := reflect.ValueOf(i)

	fmt.Println(reType, reValue)
	fmt.Printf("reType的数据类型为:%T,reValue的数据类型为:%T\n", reType, reValue)

	//如果想要获取reValue的数值，要调用Int()方法：返回v持有的有符号整数
	num2 := 80 + reValue.Int()
	fmt.Println(num2)

	//转回去
	i2 := reValue.Interface()
	//利用断言
	n := i2.(int)
	fmt.Println(n)
}
func main() {
	var num int = 100
	//将基本数据类型转到一个空接口才可以进行反射
	testReflect(num)
}
