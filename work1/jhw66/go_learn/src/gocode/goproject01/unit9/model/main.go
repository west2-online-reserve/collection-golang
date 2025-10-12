package model

import "fmt"

type Student struct {
	Name string
	Age  int
}

type person struct {
	Name string
	age  int
}

// 定义工厂模式函数，相当于构造器：
func NewPerson(name string) *person {
	p := person{name, 0}
	return &p
}

// 定义set和get方法，对age字段进行封装，因为在方法中可以加一系列限制操作，
// 确保被封装字段的安全合理性
func (p *person) SetAge(age int) {
	if age > 0 && age < 150 {
		p.age = age
	} else {
		fmt.Println("年龄范围不正确")
	}
}

func (p *person) GetAge() int {
	return p.age
}
