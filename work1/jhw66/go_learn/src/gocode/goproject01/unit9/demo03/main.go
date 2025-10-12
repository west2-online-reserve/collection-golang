package main

//方法是作用在指定的数据类型上，和指定的数据类型绑定，因此自定义类型都可以有方法，不止struct
import "fmt"

type A struct {
	Num  int
	Name string
}

//给A结构体绑定方法test
//假如方法名首字母大写，就可以在本包和其他包访问
func (a A) test() {
	fmt.Println("调用A方法test", a.Name, a.Num)
	a.Name = "HHH"
	fmt.Println("改变Name", a.Name)
}
func (a *A) test2() {
	a.Num = 7
}
func (s *A) String() string {
	str := fmt.Sprintf("Name=%v,Age=%v,调用了String方法", s.Name, s.Num)
	return str
}

//方法绑定自定义类型
type integer int

func (i integer) print() {
	fmt.Println("i=", i)
}

func main() {
	var p A = A{88, "HHH"}
	fmt.Println(p)
	p.test()
	p.Name = "hhh"
	p.test() //跟普通函数一样是值传递
	fmt.Println(p)
	p.test2()
	fmt.Println(p)

	var i integer = 20
	i.print()

	stu := A{
		Name: "qqq",
		Num:  30,
	}
	//传入地址，假如绑定了String方法,在fmt.Println就会自动调用
	fmt.Println(&stu) //注意传入的是地址
	fmt.Println(&p)
}
