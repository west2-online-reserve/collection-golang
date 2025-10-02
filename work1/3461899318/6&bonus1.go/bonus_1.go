package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

//边计算边写入文本
// func main() {
// 	file, err := os.Create("ninenine.txt") //这个函数要注意要手动关闭文件
// 	if err != nil {
// 		log.Fatal("文件创建失败", err)
// 	}
// 	defer file.Close()
// 	for i := 1; i <= 9; i++ {
// 		for j := 1; j <= i; j++ {
// 			_, err := fmt.Fprintf(file, "%d * %d = %d \t", j, i, i*j)
// 			if err != nil {
// 				log.Fatal("文件创建失败", err)
// 			}
// 		}
// 		_, err := file.WriteString("\n")
// 		if err != nil {
// 			log.Fatal("文件创建失败", err)
// 		}
// 	}
// }

// 先存入字符串中最后一次性写入
func main() {
	var str strings.Builder //注意查看定义可以看到，str是一个结构体
	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			_, err := fmt.Fprintf(&str, "%d * %d = %d \t", j, i, i*j)
			if err != nil {
				log.Fatal("文件写入失败", err)
			}
		}
		_, err := str.WriteString("\n")
		if err != nil {
			log.Fatal("文件写入失败", err)
		}
	}
	err := os.WriteFile("ninenine.txt", []byte(str.String()), 0664) //会自动关闭文件
	if err != nil {
		log.Fatal("文件写入失败", err)
	}
}
