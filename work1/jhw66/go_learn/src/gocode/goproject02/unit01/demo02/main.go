package main

import (
	"fmt"
	"os"
)

func main() {
	//os.Open 函数用于打开一个文件，并返回一个 *os.File 类型的文件对象。
	//打开文件后，我们通常需要调用 Close 方法来关闭文件，以释放系统资源。
	wd, _ := os.Getwd()
	fmt.Println("当前目录：", wd)
	file, err := os.Open("../test.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	fmt.Println("File opened successfully!")
}

// ".\test.txt"	当前目录（等同于 "test.txt")
// "..\test.txt"	上一级目录（Windows 斜杠写法）
// "D:/projects/.../unit01/test.txt"
