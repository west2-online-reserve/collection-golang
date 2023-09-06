package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	fileName := "ninenine.txt"
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()
	var mul [9][9]int
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			mul[i][j] = (i + 1) * (j + 1)
		}
	}
	for _, value := range mul {
		for i := 0; i < 9; i++ {
			strvalue := strconv.Itoa(value[i])
			fmt.Print(strvalue, " ")
			file.WriteString(strvalue)
			file.WriteString(" ")
		}
		fmt.Printf("\n")
		file.WriteString("\n")
	}
}
