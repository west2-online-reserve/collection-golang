package main

import "fmt"

type animal struct {
	Age    int
	Weight float32
}

func (an *animal) Shout() {
	fmt.Println("Shout animal!")
}
func (an *animal) ShowInfo() {
	fmt.Println(an.Age, an.Weight)
}

//结构可以使用嵌套匿名结构体所有的字段和方法
type Cat struct {
	//为了复用性，体现基础思维，嵌入匿名结构体：
	animal
}

func (c *Cat) Shout() {
	fmt.Println("Shout cat!")
}
func (c *Cat) scratch() {
	fmt.Println("scratch!")
}

type Cat2 struct {
	animal
	Cat
	string //如果不命名变量可以直接用类型
}

func (c *Cat2) String() string {
	str := fmt.Sprintf("animal里有Age:%v,Weight:%v,Cat里有animal:Age:%v,Weight:%v,string里有:%v",
		c.animal.Age, c.animal.Weight, c.Cat.Age, c.Cat.Weight, c.string)
	return str
}

func main() {
	cat := &Cat{}
	cat.animal.Age = 3
	cat.animal.Weight = 10.6
	cat.animal.Shout()
	cat.animal.ShowInfo()
	cat.scratch()

	//匿名结构体的字段访问和方法使用可以简化
	cat.Shout()        //有相同的方法时使用就近原则
	cat.animal.Shout() //不想就近原则就要这么做

	//多重继承
	// c := Cat2{animal{10, 10.1}, Cat{animal{20, 20.2}}, "miaomiaomiao"}
	c := Cat2{
		animal{
			Age:    10,
			Weight: 10.1},
		Cat{
			animal: animal{
				Age:    20,
				Weight: 20.2,
			},
		},
		"miaomiaomiao"}
	fmt.Println(c)
	fmt.Println(c.string)

	fmt.Println(&c)
}
