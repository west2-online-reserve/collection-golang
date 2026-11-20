# 使用信号量在channel里面通信

例：

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var m, n = 3, 30
	var wg sync.WaitGroup
	sem := make(chan struct{}, 1)
	cur := 1
	sem <- struct{}{}

	for i := 0; i < m; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				<-sem
				if cur > n {
					sem <- struct{}{}
					return
				} else {
					fmt.Println(cur)
					cur++
				}
				sem <- struct{}{}
			}
		}()
	}
	wg.Wait()
}

```



