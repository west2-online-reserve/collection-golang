package main

import "fmt"

func main(){
    var a,b int
    var arr[2000] int
    fmt.Scan(&a,&b)
    num:=0
    for i:=a;i<=b;i++{
        if (i%4==0&&i%100!=0)||i%400==0{
            num++
            arr[num]=i
        }
    }
    fmt.Printf("%d\n",num)
    for i:=1;i<=num;i++{
        fmt.Printf("%d ",arr[i])

    }

}
