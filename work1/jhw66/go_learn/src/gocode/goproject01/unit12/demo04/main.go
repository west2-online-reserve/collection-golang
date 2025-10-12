package main

import (
	"fmt"
	"reflect"
)

func testReflect(i interface{}) {
	reValue := reflect.ValueOf(i)
	fmt.Printf("%d,%v\n", reValue.Elem().Int(), reValue.Elem().Addr())
	reValue.Elem().SetInt(40)
}
func main() {
	var num int = 100
	testReflect(&num)
	fmt.Println(num)
}
