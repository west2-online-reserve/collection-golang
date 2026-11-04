package main

import "fmt"

func main() {
    var x, y int
    fmt.Scanf("%d %d", &x, &y)
    if !(x%4 == 0 && x%100 != 0) && !(x%400 == 0) {
    x = x + 4 - x % 4
    }
    ans := 0
    var s []int
    for i := x; i <= y; i += 4 {
        if (i % 4 == 0 && i % 100 != 0) || (i % 400 == 0) {
            s = append(s, i)
            ans++
        }
    }
    fmt.Println(ans)
    for _, v := range s {
        fmt.Print(v," ")
    }
}
