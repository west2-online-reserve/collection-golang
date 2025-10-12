// 内置函数：不需要导入包就可以使用的函数
// 内置函数存放在Builtin包下
package main

import "fmt"

func main() {
	//常见函数len new make
	fmt.Println(len("hhhh"))
	num := new(int)
	fmt.Printf("%T,%v,%v,%v", num, *num, &num, num)
}
