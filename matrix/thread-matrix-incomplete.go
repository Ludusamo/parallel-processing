package main

import (
    "fmt"
    "math/rand"
    )

const DIMX = 1000
const DIMY = 1000
var mymatrix *[DIMX][DIMY]int

func generate() {
    mymatrix=new([DIMX][DIMY]int)
    for i:=0; i< DIMX ; i++ {
        for j:=0; j<DIMY; j++ {
            mymatrix[i][j] = rand.Intn(5000)
        }
    }
    fmt.Println("Finished generating matrix")
}

func sequential_sum() {
    total := 0
    for i := 0; i< len(mymatrix); i++ {
        for j:=0; j<len(mymatrix[i]); j++ {
            total += mymatrix[i][j]
        }
    }
    fmt.Println("sequential_sum got",total)
}

func sum_row(rowSums chan int, row int) {
    sum := 0
    for _, v := range mymatrix[row] {
        sum += v
    }
    rowSums <-sum
}

func sum_rows(rowSums chan int) int {
    sum := 0
    for i := 0; i < DIMX; i++ {
        sum += <-rowSums
    }
    return sum
}

func parallel_sum() {
    // student exercise
    /*
        complete this function to print out a sum
        of the matrix. Use parallel processing
        to speed up the operation
    */
    rowSums := make(chan int)
    for i := 0; i < DIMX; i++ {
        go sum_row(rowSums, i)
    }
    total := sum_rows(rowSums)
    fmt.Println("parallel_sum got",total)
}

func main() {
    generate()
    sequential_sum()
    parallel_sum()
}
