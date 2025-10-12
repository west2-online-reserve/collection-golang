package main

import "fmt"

func main() {
	var age int
	age = 18
	fmt.Println("age=", age)

	var age2 int = 19
	fmt.Println("age2=", age2)

	//不可以使用不匹配的类型  var num int=12.56
}
