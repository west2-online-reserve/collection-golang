package main

import (
	"fmt"
	testmain "gocode/goproject01/unit5/test"
)

func test() int {
	fmt.Println("test函数被执行了!")
	return 10
}

var num int = test()

func init() {
	fmt.Println("init函数被调用了!")
}

func main() {
	fmt.Println("main函数被调用了!")
	fmt.Printf("%v %v %v", testmain.Age, testmain.Sex, testmain.Name)
	fmt.Println(num)
}
