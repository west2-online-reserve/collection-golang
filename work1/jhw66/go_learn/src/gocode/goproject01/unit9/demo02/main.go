package main

import "fmt"

type Student struct {
	Age  int
	Name string
}
type Stu Student

type Person struct {
	Age  int
	Name string
}

func main() {
	var s Student = Student{20, "nnn"}
	var p Person = Person{10, "hhh"}
	s = Student(p) //结构体是用户单独定义的类型，和其它类型进行转换时需要有完全相同的字段(名字，个数，类型)
	fmt.Println(s, " ", p)
	var S Stu = Stu{30, "mmm"}
	s = Student(S) //必须强转
	fmt.Println(s, " ", S)
}
