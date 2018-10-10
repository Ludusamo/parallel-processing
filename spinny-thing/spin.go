package main

import (
    "fmt"
    "time"
)

func fib(x int) int {
    if x <= 1 {
        return x
    }
    return fib(x - 1) + fib(x - 2)
}

func spin() {
    spinnerChars := [4]rune{ '|', '\\', '-', '/' }
    i := 0
    for {
        i = (i + 1) % 4
        fmt.Printf("%c%c", spinnerChars[i], '\r')
        time.Sleep(time.Millisecond * 100)
    }
}

func main() {
    go spin()
    fmt.Println(fib(45))
}
