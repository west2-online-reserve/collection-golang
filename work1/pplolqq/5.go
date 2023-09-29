package main
import "fmt"
func main(){
	sc:=[]int{}
	for i:=0;i<=50;i++{
		if i%3==0{
			continue
		}else{sc=append(sc,i)}
	}
	sc=append(sc,666)
	fmt.Println(sc)
}