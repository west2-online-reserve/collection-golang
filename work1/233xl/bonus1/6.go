package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Create("ninenine.txt")
	if err != nil {
		fmt.Printf("Create failed: %v", err)
		return
	}

	defer file.Close()

	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			res := fmt.Sprintf("%v*%v=%v\t", i, j, i*j)
			_, err := file.WriteString(res)
			if err != nil {
				fmt.Println("If I wont add this, VSCode would mark it yellow")
				return
			}
		}
		file.WriteString("\n")
	}

	fmt.Println("Generated successfully")
}
