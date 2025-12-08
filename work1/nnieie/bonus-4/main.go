// 通过一个 token 在通道中传递，保证打印的顺序
package main

import (
	"fmt"
	"sync"
)

// goroutine 数量
const m = 1919

// 要打印的数字总量
const n = 114514

// 共享的计数器，记录下一个要打印的数字
var counter = 1

var wg sync.WaitGroup

func main() {
	// token 用于在 goroutine 间传递信号
	token := struct{}{}

	// 创建 m 个通道，每个 goroutine 一个
	// 有 1 个缓冲，在退出时，最后一个 token 可以放进缓冲区而不会阻塞
	chs := make([]chan struct{}, m)
	for i := 0; i < m; i++ {
		chs[i] = make(chan struct{}, 1)
	}

	wg.Add(m)

	for i := 0; i < m; i++ {
		go worker(i, chs)
	}

	// 把 token 传给 goroutine 0，开始打印
	chs[0] <- token

	wg.Wait()
}

func worker(id int, chs []chan struct{}) {
	defer wg.Done()

	for {
		// 从自己的通道 chs[id] 接收 token
		token, ok := <-chs[id]
		if !ok {
			// 如果通道被关闭了，退出循环
			return
		}

		// 检查计数器是否超过 n
		if counter > n {
			// 把 token 传给下一个 goroutine
			nextID := (id + 1) % m
			chs[nextID] <- token
			// 退出循环，结束 goroutine
			return
		}

		// 打印数字
		fmt.Println(counter)
		// 更新计数器
		counter++

		// 计算下一个 goroutine 的 ID
		nextID := (id + 1) % m
		// 把 token 传给下一个 goroutine
		chs[nextID] <- token
	}
}