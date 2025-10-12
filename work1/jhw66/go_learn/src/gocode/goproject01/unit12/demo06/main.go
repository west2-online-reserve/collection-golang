package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	Name string
	Age  int
}

func TestStudentStruct(a interface{}) {
	val := reflect.ValueOf(a)
	val.Elem().Field(0).SetString("张三")
}
func main() {
	stu := Student{
		Name: "hh",
		Age:  18,
	}
	TestStudentStruct(&stu)
	fmt.Println(stu)
}
