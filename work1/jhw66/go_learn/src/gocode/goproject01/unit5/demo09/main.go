package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now()
	fmt.Println(t)
	fmt.Printf("%T\n", t)
	fmt.Printf("年为:%v\n", time.Now().Year())
	fmt.Printf("月为:%v\n", time.Now().Month())
	fmt.Printf("月为:%v\n", int(time.Now().Month()))
	fmt.Printf("日为:%v\n", time.Now().Day())
	fmt.Printf("时为:%v\n", time.Now().Hour())
	fmt.Printf("分为:%v\n", time.Now().Minute())
	fmt.Printf("秒为:%v\n", time.Now().Second())

	datest := fmt.Sprintf("当前时间为:%v-%v-%v,%d:%d:%d\n", t.Year(),
		t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	fmt.Println(datest)

	//字符串中每个数字必须固定
	datest2 := t.Format("2006/01/02 15/04/05")
	fmt.Println(datest2)
	datest3 := t.Format("2006 15:04:05")
	fmt.Println(datest3)
}
