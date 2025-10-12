package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	Name string
	Age  int
}

func (s Student) APrint() {
	fmt.Println("调用了Print()方法")
	fmt.Println("学生的名字是:", s.Name)
}

func (s Student) BGetSum(n1, n2 int) int {
	return n1 + n2
}

func (s Student) CSet(name string, age int) {
	s.Name = name
	s.Age = age
}

func TestStudentStruct(a interface{}) {
	val := reflect.ValueOf(a)
	fmt.Println(val)

	//操作结构体内部字段
	n1 := val.NumField()
	fmt.Println(n1)
	for i := 0; i < n1; i++ {
		fmt.Printf("第%d个字段的值是:%v\n", i, val.Field(i))
	}

	//操作结构体内部方法
	n2 := val.NumMethod()
	fmt.Println(n2)
	//调用方法，方法的首字母必须大写才能有对应的反射访问权限
	//方法的顺序是按照首字母的ASCII顺序排列的
	val.Method(0).Call(nil)
	//定义Value切片
	var params []reflect.Value
	params = append(params, reflect.ValueOf(10))
	params = append(params, reflect.ValueOf(20))
	result := val.Method(1).Call(params) //Call接收和返回的都是Value的切片
	fmt.Println(result[0].Int())
	params = make([]reflect.Value, 0, 2)
	params = append(params, val.Field(0))
	params = append(params, val.Field(1))
	val.Method(2).Call(params)
	fmt.Println(val)
}
func main() {
	stu := Student{
		Name: "hh",
		Age:  18,
	}
	TestStudentStruct(stu)
}
