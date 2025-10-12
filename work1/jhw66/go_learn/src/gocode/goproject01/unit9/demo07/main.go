package main

import "fmt"

//接口定义了一组方法的契约，任何类型只要实现了这些方法，就自动满足了该接口。
//接口中不能包含任何变量
type SayHello interface {
	//声明没有实现的方法
	sayHello()
}

type Chinese struct {
	name string
}
type American struct {
	name string
}

//实现接口的方法->具体实现:
//实现接口要实现所有的方法
func (person Chinese) sayHello() {
	fmt.Println("你好!")
}
func (person American) sayHello() {
	fmt.Println("hi!")
}

//中国人特有的方法
func (person Chinese) niuYangGe() {
	fmt.Println("扭秧歌")
}

//美国人特有的方法
func (person American) disce() {
	fmt.Println("disco")
}

//接受SayHello接口的函数
//定义一个函数，专门用来各国人打招呼的函数，接收具备SayHello接口能力的变量：
func greet(s SayHello) { //s为多态参数，可以通过上下文来识别具体是什么类型的实现，就体现出多态
	s.sayHello()
	//类型断言用于从接口值中提取具体类型
	//断言：可以判断是否是该类型变量 value,ok=element.(T)
	//value为变量的值 ok为bool类型 element是interface变量 T是断言类型

	//两种写法：
	// ch, flag := s.(Chinese) //判断s是否可以转换成Chinese类型并且赋给ch变量
	// if flag {
	// 	ch.niuYangGe()
	// } else {
	// fmt.Println("美国人不能扭秧歌")
	// }

	// if ch, flag := s.(Chinese); flag {
	// 	ch.niuYangGe()
	// } else {
	// 	fmt.Println("美国人不能扭秧歌")
	// }

	//多个判断：
	switch v := s.(type) { //type属于go中的一个关键字，固定写法
	case Chinese:
		v.niuYangGe()
	case American:
		v.disce()
	}

	fmt.Println("打招呼")
}

type integer int

func (i integer) sayHello() {
	fmt.Println("say hi ", i)
}

func main() {
	c := Chinese{}
	a := American{}

	greet(a)
	greet(c)

	// 将它们赋值给SayHello接口变量
	//接口本身不能创建实例，但是可以指向一个实现了该接口的自定义类型的变量
	//错误：
	// var s SayHello
	// s.sayHello()
	//正确：
	var s SayHello = c
	s.sayHello()

	//只要是自定义数据类型就可以实现接口，不仅仅是结构体类型
	var i integer = 10
	greet(i)

	//在Go语言中多态是通过接口来实现的，可以按照统一的接口来调用不同的实现，这时接口变量就呈现不同的形态
	//定义一个接口数组(多态数组),存放中国人和美国人结构体
	var arr [3]SayHello
	arr[0] = American{"00"}
	arr[1] = Chinese{"11"}
	arr[2] = Chinese{"22"}

	fmt.Println(arr)
	for i := 0; i < len(arr); i++ {
		fmt.Println(arr[i])
	}
	for _, hello := range arr {
		greet(hello)
	}
}
