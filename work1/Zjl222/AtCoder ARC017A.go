package main


import "fmt"


func isPrime(x int)bool {

	if x <= 1 {
		return false
	}

    for i:=2; i*i<=x; i++{
        if x%i == 0{
			return false
		}
	}

	return true
}

func main(){
	var x int
	fmt.Scan(&x)

	if isPrime(x)==true {
		fmt.Println("YES")
	}else {
		fmt.Println("NO")
	}
}