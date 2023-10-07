package main
import "fmt"
func main(){
	a:=0
	fmt.Scan(&a)
	if isPrime(a){fmt.Println("YES")
	}else{fmt.Println("NO")}
}
func isPrime(x int) bool{
	if x==1{return false}else if x==2{return true}else if x%2==0{return false}
	for i:=3;i<x;i+=2{
	if x%i==0{return false}}
	return true
}