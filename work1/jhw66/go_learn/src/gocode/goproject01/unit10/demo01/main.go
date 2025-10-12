package main

import (
	"fmt"
	"strconv"
	"time"
)

func test() {
	for i := 1; i < 20; i++ {
		fmt.Println("hello golong + " + strconv.Itoa(i))
		//堵塞一秒
		time.Sleep(time.Second)
	}
}
func main() { //主线程
	go test() //开启协程
	for i := 1; i < 10; i++ {
		fmt.Println("hello go + " + strconv.Itoa(i))
		//堵塞一秒
		time.Sleep(time.Second)
	}
}

//程序开始->主线程开始执行->go test()开启协程-》主线程代码继续执行->主线程结束
//                           ->执行协程中代码逻辑
//遵循主死从随
//主线程结束协程代码强制结束
