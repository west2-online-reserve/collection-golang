//5.创建一个切片(slice) 使其元素为数字1-50，从切⽚删掉数字为3的倍数的数，并且在末尾再增加⼀个数114514，输出切⽚。

package main

import "fmt"

func main() {
	s1 := []int{}
	for i := 0; i < 50; i++ {
		s1 = append(s1, i+1)
	}
	s2 := []int{}
	for j, value1 := range s1 {
		if value1%3 != 0 {
			s2 = append(s2, s1[j])
		}
	}
	s2 = append(s2, 114514)
	fmt.Print(s2)
}
