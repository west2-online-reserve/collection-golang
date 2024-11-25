package main

import "fmt"

func main(){

    var arr [10]int
    for i:=0;i<10;i++{

        fmt.Scan(&arr[i])
    }
    var t int
    fmt.Scan(&t)
    t=t+30
    sum:=0
    for i:=0;i<10;i++{
        if t>=arr[i]{
            sum++
        }
    }
    fmt.Printf("%d",sum)
}