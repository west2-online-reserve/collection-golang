package main

import "fmt"

//interface类型默认是一个指针，如果没有对interface初始化就使用，那么会输出nil
type AInterface interface {
	a()
}
type BInterface interface {
	b()
}
type CInterface interface {
	BInterface
	AInterface
	c()
}
type Stu struct {
}

func (a Stu) a() {
	fmt.Println("aaa")
}
func (b Stu) b() {
	fmt.Println("bbb")
}
func (c Stu) c() {
	fmt.Println("ccc")
}
func InterFaceA(s AInterface) {
	s.a()
}
func InterFaceB(s BInterface) {
	s.b()
}
func InterFaceC(s CInterface) {
	s.a()
	s.b()
	s.c()
}

type T interface {
}

func main() {
	var a Stu
	//一个自定义类型可以实现多个接口
	a.a()
	a.b()
	//一个接口可以继承多个别的接口，这时如果要实现这个接口，也必须将继承过来的接口全部实现
	a.c()
	InterFaceA(a)
	InterFaceB(a)
	InterFaceC(a)

	//空接口没有任何方法，所以可以把任何一个变量赋给空接口
	var E T = a
	fmt.Println(E)
	var e interface{} = 9.3
	fmt.Println(e)
}

// 如果空接口 interface{} 不包含任何方法，因此所有类型都实现了空接口。
// func PrintAnything(v interface{}) {
//     fmt.Printf("Value: %v, Type: %T\n", v, v)
// }

// func main() {
//     PrintAnything(42)           // Value: 42, Type: int
//     PrintAnything("hello")      // Value: hello, Type: string
//     PrintAnything(3.14)         // Value: 3.14, Type: float64
//     PrintAnything(Dog{Name: "Max"}) // Value: {Max}, Type: main.Dog
// }
