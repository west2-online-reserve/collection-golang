package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Create("ninenine.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			if j > 1 {
				fmt.Fprint(file, "\t")
			}
			fmt.Fprintf(file, "%d*%d=%d", j, i, j*i)
		}
		if i < 9 {
			fmt.Fprint(file, "\n")
		}
	}
}
