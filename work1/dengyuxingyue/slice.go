package main

import "fmt"

func main() {

	slice := make([]int, 0, 50)
	//make(type,size,capacity)
	for i := 1; i <= 50; i++ {
		slice = append(slice, i)
	}

	for i := 0; i < len(slice); i++ {
		if slice[i]%3 == 0 {
			slice = append(slice[:i], slice[i+1:]...)
			//通过...展开切片
			i--
		}
	}

	slice = append(slice, 114514)

	fmt.Println(slice)
}
