package main

import (
	"fmt"
	"os"
)

func main() {
	// 打开文件，如果文件不存在则创建
	file, err := os.OpenFile("ninenine.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("打开文件失败:", err)
		return
	}
	defer file.Close()

	// 生成 99 乘法表并写入文件
	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			fmt.Fprintf(file, "%d*%d=%d\t", j, i, i*j)
		}
		fmt.Fprintln(file) // 每一行结束后换行
	}

	fmt.Println("99 乘法表已生成并保存到 ninenine.txt")
}
