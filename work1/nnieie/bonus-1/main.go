package main

import (
    "fmt"
    "os"
)

func main() {
    file, err := os.Create("ninenine.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	
    for i := 1; i <= 9; i++ {
        for j := 1; j <= i; j++ {
            fmt.Fprintf(file, "%d*%d=%d\t", j, i, j*i)
        }
        fmt.Fprintln(file)
    }
}