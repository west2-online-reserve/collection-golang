package main

import (
	"fmt"
	"os"
	"log"
)

func main() {
	// 创建文件,存在则清空
	file, err := os.Create("ninenine.txt")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for i := 1; i <= 9; i++ {
		for j := i; j <= 9; j++ {
			fmt.Printf("%d*%d=%-4d", i, j, i*j)
			fmt.Fprintf(file, "%d*%d=%-4d", i, j, i*j)
		}
		fmt.Println()
		fmt.Fprintln(file, "")
	}
}
