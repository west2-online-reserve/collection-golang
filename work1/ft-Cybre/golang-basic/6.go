package main

import (
	"fmt"
	"os"
)

func main() {
	file, _ := os.Create("ninenine.txt")
	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			fmt.Fprintf(file, "%d*%d=%d", j, i, i*j)
		}
	}

}
