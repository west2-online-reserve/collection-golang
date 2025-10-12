package main

import (
	"fmt"
	"gocode/goproject01/unit9/model"
)

func main() {
	stu := model.Student{
		Name: "lll",
		Age:  99,
	}
	fmt.Println(stu)

	p := model.NewPerson("丽丽")
	p.SetAge(20)
	fmt.Println(*p)
	fmt.Println(p.GetAge())
}
