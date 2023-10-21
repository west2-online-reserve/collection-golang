package main

import "fmt"

func main() {
    var h [10]int
    var a int

    for i := 0; i < 10; i++ {
        fmt.Scan(&h[i])
    }

    fmt.Scan(&a)

    count := 0
    for i := 0; i < 10; i++ {
        if h[i] <= a+30 {
            count++
        }
    }

    fmt.Println(count)
}