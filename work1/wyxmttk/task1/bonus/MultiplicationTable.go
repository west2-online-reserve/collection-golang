package main

import (
	"fmt"
	"os"
)

func outputTable() {

	file, err := os.Create("output.txt")
	if err != nil {
		fmt.Println(err)
	}
	for i := 1; i < 10; i++ {
		for j := 1; j <= i; j++ {
			_, err := file.WriteString(fmt.Sprintf("%d * %d = %d\t", j, i, i*j))
			if err != nil {
				fmt.Println(err)
			}
			if i == j {
				_, err := file.WriteString("\n")
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}
