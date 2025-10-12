package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup //只定义不需要赋值
func main() {
	//可以在知道协程次数情况下直接wg.Add(5)
	for i := 1; i <= 5; i++ {
		wg.Add(1) //协程开始的时候加一操作
		fmt.Printf("启动协程%d\n", i)
		//使用匿名函数，直接启动协程
		go func(n int) {
			defer wg.Done() //协程结束的时候减一
			fmt.Println(n)
		}(i)
	}
	//Wait方法堵塞直到WaitGroup计时器减为0
	wg.Wait()
}
