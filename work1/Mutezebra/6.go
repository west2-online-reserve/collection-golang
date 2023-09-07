package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, _ := os.Create("ninenine.txt")
	wi := bufio.NewWriter(file)

	defer file.Close()
	defer wi.Flush()
	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			s := fmt.Sprintf("%dx%d=%d  ", j, i, j*i)
			_, _ = wi.WriteString(s)

		}
		_, _ = wi.WriteString("\n")
	}
}
