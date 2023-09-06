package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Create(".\\ninenine.txt")
	if err != nil {
		fmt.Print("Failed to Create a file")
		return
	}
	for i := 1; i < 10; i++ {
		for j := 1; j <= i; j++ {
			fmt.Fprintf(file, "%v * %v = %v  ", j, i, i*j)
		}
		fmt.Fprint(file, "\n")
	}
	file.Close()
}
