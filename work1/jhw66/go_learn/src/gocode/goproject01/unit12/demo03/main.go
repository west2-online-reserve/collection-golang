package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	Name string
	Age  int
}

func testReflect(i interface{}) {
	//调用TypeOf函数，返回reflect，Type类型数据
	reType := reflect.TypeOf(i)

	//调用ValueOf函数，返回reflect，Value类型数据
	reValue := reflect.ValueOf(i)
	fmt.Println(reType, reValue)
	fmt.Printf("reType的数据类型为:%T,reValue的数据类型为:%T\n", reType, reValue)

	//获取变量类别
	fmt.Println(reType.Kind())
	fmt.Println(reValue.Kind())
}
func main() {
	stu := Student{
		Name: "hh",
		Age:  18,
	}
	testReflect(stu)
}
