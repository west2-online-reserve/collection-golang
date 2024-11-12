package main

import "fmt"

func isPrime(x int) bool{
for i:=2;i<x/2;i++{
    if x%i==0{
        return false
    }
}
return true
}

func main(){
    var n int
    fmt.Scan(&n)
    if isPrime(n){
        fmt.Printf("YES\n")
    }else{
        fmt.Printf("NO\n")
    }

}