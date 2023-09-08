package main

import (
	"fmt"
	"time"
)

func sayHello() {
	for i := 0; i < 5; i++ {
		fmt.Println("Hello")
		time.Sleep(time.Millisecond * 100)
	}
}

func sayWorld() {
	for i := 0; i < 5; i++ {
		fmt.Println("World")
		time.Sleep(time.Millisecond * 100)
	}
}

func main() {
	go sayHello()
	go sayWorld()

	// 主协程休眠一段时间，以便让其他协程有机会执行
	time.Sleep(time.Second)

	fmt.Println("Main function exiting...")
}
