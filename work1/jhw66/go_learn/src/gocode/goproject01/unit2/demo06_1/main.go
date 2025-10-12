//GOPATH来管理文件流程(demo06_1和demo06_2来作为演示)

package main

//main包是可执行程序的入口，每个可执行程序都必须有一个main包，
//且只能有一个main包

import (
	"fmt"
	aaa_test "gocode/goproject01/unit2/demo06_2/aaa" //起别名
	"gocode/goproject01/unit2/demo06_2/test"         //导入的是包所在文件夹的路径
)

func main() {
	var n1 int = 10
	fmt.Println(n1)
	fmt.Println(test.Student) //调用时前面要定位调用的包
	fmt.Println(test.Students)
	//fmt.Println(test.teacher) //报错，teacher对外不可见
	test.Aaa()
	aaa_test.Aaa()
	//aaa.Aaa()  //报错，起别名之后原有的名字来调用函数
}

//包名：尽量保持和目录名一致
//导包：import "包路径"  包路径：从GOPATH/src开始的相对路径
//函数：func 函数名(参数列表)(返回值列表){函数体}  main函数是程序的入口函数
