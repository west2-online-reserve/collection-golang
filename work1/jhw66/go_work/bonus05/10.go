// // 思考一下m个线程打印n个数，如何保证打印的有序性
package main

import (
	"fmt"
	"sync"
)

func main() {
	n := 50
	m := 4

	var wg sync.WaitGroup
	var mu sync.Mutex
	cur := 1
	for i := 0; i < m; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				mu.Lock()
				if cur > n {
					mu.Unlock()
					return
				}
				if (cur-1)%m == id {
					fmt.Printf("num: %d\n", cur)
					cur++
				}
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
}
