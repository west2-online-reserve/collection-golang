package main

import (
	"fmt"
	"os"
)

func main() {
	f, err := os.Create("ninenine.txt")
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			//fmt.Print(j, "*", i, "=", i*j, " ")
			s := fmt.Sprintf("%d*%d=%d ", j, i, j*i)
			_, _ = f.WriteString(s)
		}
		s := fmt.Sprintf("\n")
		_, _ = f.WriteString(s)
	}

}
