package main

import (
	"fmt"
	"os"
)

func main() {
	f, err := os.Create("inenine.txt")
	if err != nil {
		fmt.Println("Create file error!")
	}
	defer f.Close()
	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			fmt.Fprintf(f, "%d * %d =%2d  ", i, j, i*j)
		}
		fmt.Fprintln(f)
	}
}
