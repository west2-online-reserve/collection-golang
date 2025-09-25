package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {

	stringBuilder := strings.Builder{}

	for x := 1; x <= 9; x++ {
		for y := 1; y <= x; y++ {
			stringBuilder.WriteString(fmt.Sprintf("%dx%d=%d ", x, y, x*y))
		}
		stringBuilder.WriteString("\n")
	}

	file, err := os.OpenFile("ninenine.txt", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error opening file", err)
		return
	}
	defer file.Close()
	_, err2 := file.WriteString(stringBuilder.String())
	if err2 != nil {
		fmt.Println("Error writing string to file:", err2)
	}
}
