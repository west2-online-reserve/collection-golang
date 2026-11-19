package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Create("ninenine.txt")
	if err != nil {
		fmt.Println("无法创建文件:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			equation := fmt.Sprintf("%d×%d=%-2d  ", j, i, i*j)
			_, err := writer.WriteString(equation)
			if err != nil {
				fmt.Println("失败:", err)
				return
			}
		}
		_, err := writer.WriteString("\n")
		if err != nil {
			fmt.Println("失败:", err)
			return
		}
	}

	err = writer.Flush()
	if err != nil {
		fmt.Println("失败:", err)
		return
	}

	fmt.Println("99乘法表已保存到ninenine.txt文件")
}