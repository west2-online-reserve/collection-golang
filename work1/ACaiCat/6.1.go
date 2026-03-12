package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.OpenFile("ninenine.txt", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error opening file", err)
		return
	}

	defer file.Close()

	writer := bufio.NewWriter(file)
	defer func(writer *bufio.Writer) {
		err := writer.Flush()
		if err != nil {
			fmt.Println("Error flushing file:", err)
		}
	}(writer)

	for x := 1; x <= 9; x++ {
		for y := 1; y <= x; y++ {
			_, err = fmt.Fprintf(writer, "%dx%d=%d ", x, y, x*y)
			if err != nil {
				fmt.Println("Error writing to buffer:", err)
			}
		}
		_, err := writer.WriteString("\n")
		if err != nil {
			fmt.Println("Error writing to buffer:", err)
		}
	}

}
