package main

import (
	"fmt"
	"log"
	"os"
)

// 文件信息操作
func main() {
	fileInfo, err := os.Stat("../test.txt")
	if os.IsNotExist(err) {
		log.Fatal("文件不存在")
	} else {
		fmt.Println("文件存在")
	}

	fmt.Println("文件名:", fileInfo.Name())
	fmt.Println("文件大小:", fileInfo.Size(), "字节")
	fmt.Println("权限:", fileInfo.Mode())
	fmt.Println("最后修改时间:", fileInfo.ModTime())
	fmt.Println("是目录吗:", fileInfo.IsDir())

	err = os.Rename("../test.txt", "new.txt")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("重命名成功")
	err = os.Rename("new.txt", "../test.txt")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("重命名成功")
}
