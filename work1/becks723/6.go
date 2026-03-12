package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Create("ninenine.txt")
	if err != nil {
		fmt.Println("Error creaing file: ", err)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "   %3d%3d%3d%3d%3d%3d%3d%3d%3d\n", 1, 2, 3, 4, 5, 6, 7, 8, 9)
	for opcode := 1; opcode <= 9; opcode++ {
		fmt.Fprintf(file, "%3d", opcode)
		for i := 0; i < opcode; i++ {
			fmt.Fprintf(file, "%3d", (i+1)*opcode)
		}
		fmt.Fprintln(file) // line break
	}
}
