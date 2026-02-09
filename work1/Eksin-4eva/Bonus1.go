package main

import (
	"fmt"
	"os"
)

func main() {
	const fileName = "ninenine.txt"
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error:%v", err)
		return
	}

	defer file.Close()

	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			fmt.Fprintf(file, "%d*%d=%-2d ", j, i, i*j)
		}
		file.WriteString("\n")
	}

}
