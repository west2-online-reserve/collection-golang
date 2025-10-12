package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup
var lock sync.RWMutex

// 读写锁在只有读的时候锁不生效 有读有写的时候才生效
func read() {
	defer wg.Done()
	lock.RLock()
	fmt.Println("开始读取数据")
	time.Sleep(time.Second)
	fmt.Println("读取数据成功")
	lock.RUnlock()
}

func write() {
	defer wg.Done()
	lock.Lock()
	fmt.Println("开始修改数据")
	time.Sleep(time.Second)
	fmt.Println("修改数据成功")
	lock.Unlock()
}

func main() {
	wg.Add(6)
	for i := 0; i < 5; i++ {
		go read()
	}
	go write()
	wg.Wait()
}
