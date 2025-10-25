package main

import "fmt"

//
//Go语言中的切片和数组的区别有哪些？
/*------------------------------------------------------------------*/
//数组：长度是固定的，在声明时必须指定长度，且后续无法修改。
//数组：类型由元素类型 + 长度共同决定。例如[3]int和[5]int是两种不同的数组类型，不能直接赋值或比较。
//数组：是值类型。当数组被赋值、作为函数参数传递时，会复制整个数组的所有元素。若数组体积较大，会导致内存开销增加和性能损耗。
//
/*------------------------------------------------------------------*/
//切片：长度是动态的，声明时无需指定长度（或通过make指定初始长度），可以通过append等操作动态改变长度。
//切片：类型仅由元素类型决定，与长度无关。例如[]int是一种切片类型，无论长度是 3 还是 5，类型都相同，可直接赋值。
//切片：是引用类型。切片本身不存储数据，它包含一个指向底层数组的指针、长度（len）和容量（cap）。当
//切片被赋值或传递时，仅复制这 3 个字段（指针、len、cap），不会复制底层数据，因此效率更高。
/*------------------------------------------------------------------*/

// 测试所以切片和map的初始化方式
func main() {
	var slice1 []int //长度容量都为0

	var slice2 []int = []int{1, 2, 3, 4} //长度和容量与初始化元素数量相同

	var slice5 []int = make([]int, 4, 4)

	slice3 := []int{1, 2, 3, 4, 5} //长度和容量与初始化元素数量相同
	//slice可为slice2的简化写法

	slice4 := make([]int, 3, 3) //使用make初始化

	slice6 := slice3[1:3] //这种初始化的长度与容量与其他有差异
	//当然还有这样的初始化方法
	//slice6 := slice3[:3]
	//slice6 := slice3[0:]

	//切片统一测试
	fmt.Println(slice1)
	fmt.Println(len(slice2), cap(slice2), slice2)
	fmt.Println(len(slice3), cap(slice3), slice3)
	fmt.Println(len(slice4), cap(slice4), slice4)
	fmt.Println(len(slice5), cap(slice5), slice5)
	fmt.Println(len(slice6), cap(slice6), slice6)
	fmt.Println("\n")

	var m1 map[string]int = map[string]int{"1243": 1}
	m2 := map[string]int{"1243": 1, "ad": 2} //可看成m1的简化初始化方法
	var m3 map[string]int                    //这种初始化是不能直接赋值的
	var m4 = make(map[string]int, 10)        //可以设置容量
	//m3["312"] = 2 这里会panic
	m1["ad"] = 2
	m4["adc"] = 5
	fmt.Println(m1)
	fmt.Println(m2)
	fmt.Println(m3)
	fmt.Println(m4)
}
