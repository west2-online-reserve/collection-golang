package main

import (
	"fmt"
	"os"
)

func main() {
	file, _ := os.Create("ninenine.txt")
	defer file.Close()

	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			result := i * j
			fmt.Fprintf(file, "%d*%d=%d ", i, j, result)
		}
		fmt.Fprintln(file)
	}

}

