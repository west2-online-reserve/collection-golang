package main

import "fmt"

//结构体定义
type Teacher struct {
	Name   string
	Age    int
	School string
}

func main() {
	//创建对象
	var ma Teacher
	fmt.Println(ma)
	ma.Name = "ma"
	ma.Age = 22
	ma.School = "qinghua"
	fmt.Println(ma)

	var ha Teacher = Teacher{"ha", 23, "beida"}
	fmt.Println(ha)

	var t *Teacher = new(Teacher)
	(*t).Name = "sa"
	(*t).Age = 45
	t.School = "zhejiangu" //底层转化为(*t).School
	fmt.Println(*t)

	var r *Teacher = &Teacher{"rr", 11, "fudau"}
	fmt.Println(r, " ", *r)

	var l Teacher = Teacher{
		Name: "ll",
		Age:  33,
	}
	fmt.Println(l)
}
