package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup
var lock sync.RWMutex

func writeData(intChan chan int) {
	defer wg.Done()
	lock.RLock()
	for i := 1; i <= 10; i++ {
		intChan <- i
		fmt.Println("写入的数据为：", i)
	}
	lock.RUnlock()
	close(intChan)
}
func readData(intChan chan int) {
	defer wg.Done()
	for v := range intChan {
		fmt.Println("读取的数据为：", v)
	}
}

func main() {
	intChan := make(chan int, 100)

	wg.Add(2)
	go writeData(intChan)
	go readData(intChan)
	fmt.Printf("hhh!\n")

	wg.Wait()
}
