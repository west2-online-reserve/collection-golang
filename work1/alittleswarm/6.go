package main

import (
	"fmt"
	"os"
)

func main() {
	var f *os.File
	f, _ = os.OpenFile("./ninenine.txt", os.O_CREATE|os.O_WRONLY, 0666)
	for i := 1; i < 10; i++ {
		for j := 1; j <= i; j++ {
			s := fmt.Sprintf("%d*%d=%d\t", j, i, i*j)
			print(s)
			f.WriteString(s)
		}
		f.Write([]byte("\n"))
		println()
	}
}
