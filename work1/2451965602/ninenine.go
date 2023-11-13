package main
import (
	"fmt"
	"os"
)
func main(){
	file,_ := os.Create("ninenine.txt")
	for x := 1; x <= 9; x++{
		for y := 1; y <= x; y++{
			fmt.Fprintf(file,"%d * %d = %d ",y,x,x*y)
		}
		fmt.Fprintln(file,"\n")
	}
	defer file.Close()
}
