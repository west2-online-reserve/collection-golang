package main

import "fmt"

func main() {
    ans := 0
    var height [10]int
    for i := 0; i < 10; i++ {
        fmt.Scan(&height[i])
    }
    var tt int
    fmt.Scan(&tt)
    for i := 0; i < 10; i++ {
        if (tt + 30) >= height[i] {
            ans++
        }
    }
    fmt.Println(ans)
}
