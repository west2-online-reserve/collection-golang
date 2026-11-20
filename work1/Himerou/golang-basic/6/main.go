package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Create("ninenine.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	for i := 1; i < 10; i++ {
		for j := 1; j <= i; j++ {
			_, err := file.WriteString(fmt.Sprintf("%d * %d = %d ", i, j, i*j))
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		_, err = file.WriteString("\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}
