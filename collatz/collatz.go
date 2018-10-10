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

func main() {
	maxN, maxLen := 1, 1
	for i := 1; i < 10000000; i++ {
		length := collatzLen(i)
		if length > maxLen {
			maxN, maxLen = i, length
		}
	}
	fmt.Printf("Longest sequence starts at %d, length %d\n", maxN, maxLen)
}
