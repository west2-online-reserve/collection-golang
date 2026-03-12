package main

import (
	"fmt"
)

const N = 1000000 + 1 //const当var用，所以不用：

func main() {
	fmt.Println("请输入数组的元素数量：")
	var n int
	fmt.Scanln(&n)        //scanln不吃'\n'如果还有输入，要先从处理换行
	s := make([][]int, N) //二维切片至多n行
	for i := 0; i < n; i++ {
		x := 0
		fmt.Scanf("%d", &x)    //我们不关心数组的值是多少，没有保存的必要
		s[x] = append(s[x], i) //做预处理
	}
	target := 0
	fmt.Scan(&target)
	for i := 0; i <= target/2; i++ {
		len1 := len(s[i])
		len2 := len(s[target-i])
		if target-i == i {
			for j := 0; j < len1; j++ {
				for m := j + 1; m < len2; m++ {
					fmt.Printf("[%d,%d]", s[i][j], s[target-i][m])
				}
			}
		} else if len1 > 0 && len2 > 0 {
			for j := 0; j < len1; j++ {
				for m := 0; m < len2; m++ {
					fmt.Printf("[%d,%d]", s[i][j], s[target-i][m])
				}
			}
		}
	}
}
