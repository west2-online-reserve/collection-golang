package main
import "fmt"
func main(){
	a,b:=0,0
	fmt.Scan(&a,&b)
	sc:=[]int{}
	nums:=0
	for i:=a;i<=b;i++{
		if i%4==0 && i%100!=0{
			sc=append(sc,i)}
		if i%100==0 && i%400==0{sc=append(sc,i)}}
	fmt.Println(nums)
	for i:=range sc{
		fmt.Printf("%d ",sc[i])}
}