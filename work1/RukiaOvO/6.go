package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	var temp string

	file, err := os.OpenFile("./ninenine.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			temp = strconv.Itoa(j) + "*" + strconv.Itoa(i) + "=" + strconv.Itoa(i*j) + " "
			file.WriteString(temp)
		}
		if i != 9 {
			file.WriteString("\n")
		}
	}

	defer file.Close()
}
