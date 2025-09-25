package main

import (
	"fmt"
	"sync"
)

var (
	number int
	mu     sync.Mutex
)

func main() {
	m := 6
	n := 114514
	ch := make(chan struct{}) // 创建一个管道用来等待打印完成
	for i := 0; i < m; i++ {
		go printNumber(ch, n)
	}
	<-ch // 未没有打印完成时阻塞主程序

}
func printNumber(ch chan struct{}, n int) {
	for {
		mu.Lock()
		if number == n {
			ch <- struct{}{} // 向管道发送空结构体以结束主程序
			mu.Unlock()      // 打开锁以便其他goroutine释放 (感觉是不是有点咸鱼了)
			return
		}
		number++
		fmt.Println(number)
		mu.Unlock()
	}
}
