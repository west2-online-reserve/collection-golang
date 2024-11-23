package main

import (
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

	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			line := fmt.Sprintf("%d * %d = %d\t", j, i, j*i)
			_, err := file.WriteString(line)
			if err != nil {
				fmt.Println("写入文件时出错:", err)
				return
			}
		}
		_, err := file.WriteString("\n")
		if err != nil {
			fmt.Println("写入文件时出错:", err)
			return
		}
	}
}
