package testmain

import "fmt"

var Age int
var Sex string
var Name string

func init() {
	fmt.Println("test中的init函数被执行了!")
	Age = 19
	Sex = "nnn"
	Name = "mmm"
}
