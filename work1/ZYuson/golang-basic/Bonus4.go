// 使用循环传递令牌的方式保证多线程打印的有序性

package main

import (
	"fmt"
	"sync"
)

func main() {
	const m = 3
	const n = 10

	ch := make([]chan struct{}, m)	// 创建m个channel（空结构体类型channel，节省内存）
	for i := range ch {
		ch[i] = make(chan struct{}, 1) // 每个channel分配容量1
	}

	var wg sync.WaitGroup	// 创建WaitGroup
	for id := 0; id < m; id++ {
		wg.Add(1)	// 增加WaitGroup计数
		go func(id int) {
			defer wg.Done()	//减少WaitGroup计数
			for x := id + 1; x <= n; x += m {
				<-ch[id]                       // 等自己的令牌
				fmt.Println(x)                 // 拿到令牌后打印
				ch[(id+1)%m] <- struct{}{}     // 令牌传给下一个channel（空结构体作令牌）
			}
		}(id)
	}

	ch[0] <- struct{}{} // 从0号channel开始传递令牌
	wg.Wait()	// 等待所有goroutine完成
}