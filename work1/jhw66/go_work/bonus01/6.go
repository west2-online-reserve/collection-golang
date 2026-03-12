//写一个99乘法表，并且把结果保存到同⽬录下ninenine.txt，⽂件保存命名为"6.go"。

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.OpenFile("ninenine.txt", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			writer.WriteString(fmt.Sprintf("%d+%d=%d\t", j, i, i+j))
		}
		writer.WriteString("\n")
	}
	writer.Flush()

}
