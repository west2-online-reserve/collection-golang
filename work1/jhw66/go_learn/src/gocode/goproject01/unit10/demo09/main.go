package main

import (
	"fmt"
	"time"
)

func printNum() {
	for i := 0; i <= 10; i++ {
		fmt.Println(i)
	}
}
func devide() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	num1 := 10
	num2 := 0
	result := num1 / num2
	fmt.Println(result)
}

func main() {
	go printNum()
	go devide()

	time.Sleep(time.Second * 3)
}
