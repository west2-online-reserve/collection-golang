package main

import "fmt"

func isPrime (x int) bool {
	is:=1
	for i:=2;i<x/2;i++ {
		if x%i==0 {
			is=0
			break
		}
	}
	if is==1 {
		return true
	}else{
		return false
	}
}
func main(){
	var n int
	fmt.Scan(&n)
	if isPrime(n) {
		fmt.Printf("YES\n")
	}else{
		fmt.Printf("NO\n")
	}

}