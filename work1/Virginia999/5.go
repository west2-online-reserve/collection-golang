// 向切片中增删元素
package main

import "fmt"

func main() {
	a := []int{}
	start := 1
	end := 50
	for i := start; i <= end; i++ {
		a = append(a, i)
	}
	//易错点1：第一次循环访问i=1会触发索引越界   （b := a[i]）
	//易错点2：从后往前索引，避免索引错位
	for i := len(a) - 1; i >= 0; i-- {
		//易错点3：是i对应下标表示的a，不是i本身，否则同样造成索引错位
		if a[i]%3 == 0 {
			a = append(a[:i], a[i+1:]...)
		}
	}
	b := 114514
	a = append(a, b)
	fmt.Println(a)
}
