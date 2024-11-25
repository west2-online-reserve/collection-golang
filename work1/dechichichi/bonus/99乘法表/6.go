package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Create("ninenine.txt")
	if err == nil {
		for i := 1; i <= 9; i++ {
			for j := 1; j <= 9; j++ {
				fmt.Fprintf(file, "%d*%d=%d\t", i, j, i*j)
			}
			fmt.Fprintf(file, "\n")
		}
	}
	file.Close()
}
