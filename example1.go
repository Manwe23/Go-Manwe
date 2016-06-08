package main

import "fmt"

//some simple fibonacci function
func fibonacci() func() int {
    a,b := 0,1
    return func() int {
        tmp := a
        a = b
        b = tmp + b
        return b
    }
}

func main() {
    f := fibonacci()
    for i := 0; i < 10; i++ {
        fmt.Println(f())
    }
}
