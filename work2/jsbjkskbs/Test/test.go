package main

import (
	"log"
	"net/http"
	"time"
)

// GoLimit(协程限制)
// ----------------------------------------------------
type GoLimit struct {
	cnt chan int
}

func NewGoLimit(maxRoutine int) *GoLimit {
	return &GoLimit{make(chan int, maxRoutine)}
}

func (goLimit *GoLimit) Add() {
	goLimit.cnt <- 1
}

func (goLimit *GoLimit) Done() {
	<-goLimit.cnt
}
//----------------------------------------------------

func main() {
	rLimit := NewGoLimit(512)
	reqCnt := 0
	go func() {
		for {
			rLimit.Add()
			go func() {
				defer rLimit.Done()
				for {
					reqCnt++
					http.Get("https://api.bilibili.com/x/v2/reply/main?oid=318291485&type=1&mode=2&next=1000")
				}
			}()
		}
	}()
	var latestCnt int
	go func() {
		for {
			cntInGap := reqCnt - latestCnt
			latestCnt = reqCnt
			log.Printf("cnt in one second: %d \n", cntInGap)
			time.Sleep(time.Second)
		}
	}()
	for {}
}
