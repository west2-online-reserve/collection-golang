package main
import "fmt"
func main(){
	var slice []int
	slice=make([]int,50,100)
	for i:=1;i<=50;i++{
		slice= append(slice, i)
	}
	fmt.Println(slice)
}
