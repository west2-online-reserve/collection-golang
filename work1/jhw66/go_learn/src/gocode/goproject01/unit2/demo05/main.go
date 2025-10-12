package main

import "fmt"

func main() {
	var age int = 19
	var ptr *int = &age
	fmt.Println(&age)
	fmt.Println(ptr)
	fmt.Println(*ptr)
	fmt.Println(&ptr)

	*ptr = 20
	fmt.Println(age)

}
