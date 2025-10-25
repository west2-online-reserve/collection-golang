package main

import (
	"fmt"
	"sync"
)

func relayPrint(n, m int) {
	chs := make([]chan struct{}, m)
	for i := range chs {
		chs[i] = make(chan struct{}, 1)
	}

	var wg sync.WaitGroup
	wg.Add(m)

	for i := 0; i < m; i++ {
		go func(id int) {
			defer wg.Done()
			num := id + 1

			for {
				_, ok := <-chs[id]
				if !ok {
					return
				}

				if num > n {
					for j := 0; j < m; j++ {
						if j != id {
							close(chs[j])
						}
					}
					return
				}

				fmt.Printf("%d ", num)
				num += m

				nextID := (id + 1) % m
				chs[nextID] <- struct{}{}
			}
		}(i)
	}

	chs[0] <- struct{}{}
	wg.Wait()
}

func main() {
	relayPrint(10, 3)
}
