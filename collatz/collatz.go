package main

import (
	"fmt"
)

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

/** Generates numbers from i to j (non-inclusive) fed through a channel
 * @param i start of the sequence
 * @param j end of the sequence (non-inclusive)
 * @return read only channel of the sequence
 */
func generateNumbers(i int, j int) <-chan int {
	seqChan := make(chan int)
	go (func() {
		for n := i; n < j; n++ {
			seqChan <- n
		}
	})()
	return seqChan
}

func main() {
	numChan := generateNumbers(1, 10000001)
	maxN, maxLen := 1, 1
	for i := 1; i < 10000001; i++ {
		n := <-numChan
		length := collatzLen(n)
		if length > maxLen {
			maxN, maxLen = n, length
		}
	}
	fmt.Printf("Longest sequence starts at %d, length %d\n", maxN, maxLen)
}
