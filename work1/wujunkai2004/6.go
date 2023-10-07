package main

import "fmt"
import "os"

func main(){
	file, err := os.Create("ninenine.txt");
	if err != nil {
		fmt.Printf("File IO Error: %s", err);
		return;
	}
	defer file.Close()

	for one:=1; one<=9; one++ {
		for two:=one; two<=9; two++ {
			fmt.Fprintf(file, "%dx%d=%2d\t", one, two, one*two )
		}
		fmt.Fprintf(file, "\r\n")
	}
}