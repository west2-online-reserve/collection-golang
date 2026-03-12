//1.写一个99乘法表，并且把结果保存到同⽬录下ninenine.txt，⽂件保存命名为"6.go"。

package main

import (
	"fmt"
	"os"
)

func main() {

	file, err := os.Create("D:\\Golang_code\\Lab1\\ninenine.txt")
	if err != nil {
		fmt.Println("文件创建失败", err)
		return
	}
	defer file.Close()

	for i := 1; i <= 9; i++ {

		for j := 1; j <= i; j++ {
			content := fmt.Sprintf("%d*%d=%d ", i, j, i*j)
			_, err = file.WriteString(content)
			if err != nil {
				fmt.Println("文本输入失败：", err)
				return
			}
		}

		_, err := file.WriteString("\n")
		if err != nil {
			fmt.Println("文本输入失败：", err)
			return
		}
	}

}
