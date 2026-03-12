//可以通过mutex，在同时启用多个线程时确保只有一个线程发挥作用，从而确保打印的有序性

package main

import (
	"fmt"
	"sync"
)

var locker sync.Mutex
var count int = 0
var wg sync.WaitGroup

func PrintNumbers(n int) {
	for i := 0; i < n; i++ {
		locker.Lock()
		wg.Add(1)
		count++
		fmt.Printf("%d\t", count)
		wg.Done()
		locker.Unlock()

	}
}
func main() {
	go PrintNumbers(10)
	go PrintNumbers(10)
	go PrintNumbers(10)
	defer wg.Wait()
}
