package main

import (
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Create("ninenine.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var builder strings.Builder
	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			builder.WriteString(strconv.Itoa(j))
			builder.WriteByte('*')
			builder.WriteString(strconv.Itoa(i))
			builder.WriteByte('=')
			builder.WriteString(strconv.Itoa(i * j))
			builder.WriteByte(' ')
		}
		builder.WriteByte('\n')
	}
	_, err = file.WriteString(builder.String())
	if err != nil {
		panic(err)
	}
}
