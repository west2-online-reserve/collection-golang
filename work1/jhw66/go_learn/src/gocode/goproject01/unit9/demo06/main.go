package main

import "fmt"

type A struct {
	a int
	b string
}

type B struct {
	c int
	d string
	a int
}

type C struct {
	*A
	B
	int
}
type D struct {
	a int
	b string
	c B //组合模式
}

func main() {
	c := C{&A{10, "aaa"}, B{20, "bbb", 50}, 999}
	fmt.Println(c)
	fmt.Println(*c.A)

	d := D{10, "ooo", B{66, "ppp", 99}}
	fmt.Println(d)

}
