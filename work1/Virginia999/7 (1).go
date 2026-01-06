// 编写一个判断质数的函数，并在主函数中引用它
package main

import "fmt"

func Select(a int) bool {
	result := false
	for i := 0; i < a; i++ {
		for j := 0; j < a; j++ {
			if a == i*j {
			} else if a != i*j {
				result = true
			}
		}
	}
	return result
}
func main() {
	var a int
	fmt.Scanf("%d", &a)
	select1 := Select(a)
	fmt.Println(select1)
}
