package main

import (
	"fmt"
	"os"
)

func write() {
	fileObj, err := os.OpenFile("ninenine.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open file failed")
		return
	}
	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			str := fmt.Sprintf("%d*%d=%d ", j, i, i*j)
			fileObj.WriteString(str)
		}
		str := "\n"
		fileObj.WriteString(str)
	}
}

func main() {
	write()
}
