package main

import (
	"fmt"
	"os"
)

func main() {
	path, _ := os.Getwd()
	path = path + "/ninenine.txt"
	file, _ := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	defer file.Close()

	for i := 1; i <= 9; i++ {
		for j := 1; j <= 9; j++ {
			str := fmt.Sprintf("%d*%d=%d\n", i, j, i*j)
			file.WriteString(str)
		}
	}
}
