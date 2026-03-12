// 写一个99乘法表，并且把结果保存到同⽬录下ninenine.txt，⽂件保存命名为"6.go"。

package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {

	var str strings.Builder
	for i := 1; i <= 9; i++ {
		for j := i; j <= 9; j++ {
			str.WriteString(fmt.Sprintf("%v*%v=%v\t", i, j, i*j))
		}
		str.WriteString("\n")
	}

	os.WriteFile("ninenine.txt", []byte(str.String()), 0666)
}
