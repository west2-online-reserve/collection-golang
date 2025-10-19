package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("./ninenine.txt")
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 1; i <= 10; i++ {
		for j := 1; j <= 10; j++ {
			str := fmt.Sprintf("%d * %d = %d", i, j, i*j)
			_, err := file.WriteString(str)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

}
