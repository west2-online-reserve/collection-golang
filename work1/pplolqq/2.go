package main
import "fmt"
func main(){
	sc:=[]int{}
	a:=0
	for i:=0;i<8;i++{//输入8组数据
		fmt.Scan(&a)
		sc=append(sc,a)
	}
	p:=0
	fmt.Scan(&p)
	p+=30
	var nums int
	for i:=0;i<8;i++{
		if sc[i]<=p{nums++}
	}
	fmt.Println(nums)
}