package main

import (
	"fmt"
	"os"
)

func main() {
	var res string

	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			res += fmt.Sprintf("\t%d*%d=%d", i, j, i*j)
		}

		res += "\n"
	}

	os.WriteFile("ninenine.txt", []byte(res), 0644)
}
