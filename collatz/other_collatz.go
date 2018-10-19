package main

import (
	"fmt"
	"time"
)

type WrapperFunc func()

/** Takes a function and times it and reports how long it took to run
 * @param id string identifier for printing purposes
 * @param f function to be run
 */
func timeit(id string, f WrapperFunc) {
	start := time.Now()
	f()
	elapsed := time.Since(start)
	fmt.Printf("%s ran in %s\n", id, elapsed)
}

type CollatzPair struct {
	start,
	length int
}

/** Calculates the Collatz Length of an integer n recursively
 * @param n the starting number
 * @return the Collatz Length of n
 */
func collatzLen(n int) int {
	if n <= 0 {
		return -1
	}
	if n == 1 {
		return 1
	}
	if n%2 == 0 {
		return collatzLen(n/2) + 1
	}
	return collatzLen(n*3+1) + 1
}

func maxCollatz(start int, end int, step int, resultChan chan CollatzPair) {
	max := CollatzPair{0, 0}
	for n := start; n < end; n += step {
		l := collatzLen(n)
		if l > max.length {
			max = CollatzPair{n, l}
		}
	}
	resultChan <- max
}

func main() {
	timeit("other", func() {
		numProc := 6
		resultChan := make(chan CollatzPair, numProc)
		for i := 0; i < numProc; i++ {
			go maxCollatz(i, 10000001, numProc, resultChan)
		}
		max := CollatzPair{0, 0}
		for i := 0; i < numProc; i++ {
			pair := <-resultChan
			if max.length < pair.length {
				max = pair
			}
		}
		fmt.Printf("Longest sequence starts at %d, length %d\n", max.start, max.length)
	})
}
