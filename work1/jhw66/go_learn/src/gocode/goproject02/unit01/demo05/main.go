// 文件复制
package main

import (
	"io"
	"log"
	"os"
)

func main() {
	srcFile, err := os.Open("../test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create("destination.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	bytesCopied, err := io.Copy(dstFile, srcFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("复制完成，共复制 %d 字节", bytesCopied)
}
