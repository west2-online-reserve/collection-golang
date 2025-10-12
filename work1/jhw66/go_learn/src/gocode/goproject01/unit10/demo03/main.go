package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup
var TotalNum int = 0

// 加入互斥锁
var lock sync.Mutex

func add() {
	defer wg.Done()
	lock.Lock()
	for i := 0; i < 100000; i++ {
		TotalNum++
	}
	fmt.Println("Overadd")
	fmt.Printf("Overadd:%d\n", TotalNum) //如果不加lock结果不会一定是100000因为是并发进行
	lock.Unlock()
}
func sub() {
	defer wg.Done()
	for i := 0; i < 100000; i++ {
		lock.Lock()
		TotalNum--
		lock.Unlock()
	}
	fmt.Println("OverDnoe")
	fmt.Printf("OverDone:%d\n", TotalNum)

}
func main() {
	wg.Add(2)
	//协程间是并发执行，不加锁会导致取值混乱
	go add()
	//time.Sleep(time.Second * 2)
	fmt.Println(TotalNum) //这个才一定是100000因为是主线程休息了2s，add协程已经结束
	go sub()
	wg.Wait()
	fmt.Println(TotalNum)
}
