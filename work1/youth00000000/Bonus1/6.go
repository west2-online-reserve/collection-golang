package main

import (
	"log"
	"os"
	"strconv"
)

func main() {
	write()
}

func write() {
	file, err := os.Create("ninenine.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for i := 1; i < 10; i++ {
		for j := 1; j <= i; j++ {
			_, err := file.WriteString(strconv.Itoa(i) + "*" + strconv.Itoa(j) + "=" + strconv.Itoa(i*j) + "\t")
			if err != nil {
				log.Fatal(err)
			}
		}
		file.WriteString("\n")
	}
}
