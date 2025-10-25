/*
bonus5.go -- m个线程按顺序打印n个数
*/

package main

import (
	"fmt"
	"sync"
)

func main() {
	m := 3
	n := 100

	var wg sync.WaitGroup
	chans := make([]chan struct{}, m)
	for i := 0; i < m; i++ {
		chans[i] = make(chan struct{})
	}

	for i := 0; i < m; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for k := id + 1; k <= n; k += m {
				<-chans[id]
				fmt.Printf("Routine %d prints %d\n", id, k)
				if k < n {
					chans[(id+1)%m] <- struct{}{}
				} else { // k == n
					for _, ch := range chans {
						close(ch)
					}
				}

			}
		}(i)
	}

	chans[0] <- struct{}{}
	wg.Wait()
}
