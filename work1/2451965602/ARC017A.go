package main
import "fmt"
var x int
var result string
func isPrime(x int) bool{
	for i := 2; i < x ; i++{
		if x % i == 0 {
			return false
		}
	}
	return true
}
func main(){
	fmt.Scan(&x)
	ret := isPrime(x)
	if ret{
		fmt.Println("YES")
	}else{
		fmt.Println("NO")
	}
}