package main

import "fmt"

func main() {

	var n1 int = +10
	fmt.Println(n1)
	var n2 int = 4 + 7
	fmt.Println(n2)
	var s1 string = "hello" + "world"
	fmt.Println(s1)

	fmt.Println(10 / 3)   //整数除法结果为整数
	fmt.Println(10.0 / 3) //浮点数除法结果为浮点数

	fmt.Println(10 % 3)
	fmt.Println(-10 % 3)
	fmt.Println(10 % -3)
	fmt.Println(-10 % -3)

	n1++
	fmt.Println(n1)
	n1--
	fmt.Println(n1)

	fmt.Println(5 == 9)
	fmt.Println(5 != 9)
	fmt.Println(5 > 9)
	fmt.Println(5 < 9)
	fmt.Println(5 >= 9)
	fmt.Println(5 <= 9)

	fmt.Println(true && false)
	fmt.Println(true || false)
	fmt.Println(!true)
}
