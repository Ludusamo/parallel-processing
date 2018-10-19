package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

var cache = make(map[int]int)
var cacheMutex = &sync.Mutex{}

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
	start  int
	length int
}

/** Calculates the Collatz Length of an integer n
 * @param n the starting number
 * @return the Collatz Length of n
 */
func collatzLen(n int) int {
	cur := n
	encountered := make([]int, 0)
	length := 1
	for cur != 1 {
		cacheMutex.Lock()
		l, has := cache[cur]
		cacheMutex.Unlock()
		if has {
			length += l - 1
			break
		}
		encountered = append(encountered, cur)
		length++
		if cur%2 == 0 {
			cur /= 2
		} else {
			cur = cur*3 + 1
		}
	}
	cacheMutex.Lock()
	for i, v := range encountered {
		cache[v] = length - i
	}
	cacheMutex.Unlock()
	return length
}

/** Takes input from an input channel and places Collatz Pair into output chan
 * @param inChan the input channel
 * @param outChan the output channel
 */
func collatzLenProducer(inChan <-chan int, outChan chan *CollatzPair) {
	for {
		n := <-inChan
		outChan <- &CollatzPair{n, collatzLen(n)}
	}
}

/** Generates numbers from i to j (non-inclusive) fed through a channel
 * @param i start of the sequence
 * @param j end of the sequence (non-inclusive)
 * @return read only channel of the sequence
 */
func generateNumbers(i int, j int) <-chan int {
	seqChan := make(chan int, 5)
	go (func() {
		for n := i; n < j; n++ {
			seqChan <- n
		}
	})()
	return seqChan
}

/** Finds the maximum of a number of values coming from a channel
 * @param numIncoming the number of values to expect
 * @param inChan the input channel
 * @return the maximum Collatz Pair flowing through the channel
 */
func maxOfChan(numIncoming int, inChan chan *CollatzPair) *CollatzPair {
	maxPair := &CollatzPair{0, 0}
	for i := 0; i < numIncoming; i++ {
		pair := <-inChan
		if pair.length > maxPair.length {
			maxPair = pair
		}
	}
	return maxPair
}

/** Finds maximum Collatz Length between two numbers in a parallel architecture
 * @param start the beginning of the sequence to check
 * @param end the end of the sequence non-inclusive
 * @param numWorkers number of worker threads to spawn
 * @return the maximum Collatz number and its length
 */
func maxCollatzParallel(start int, end int, numWorkers int) *CollatzPair {
	numGen := generateNumbers(start, end)
	lenChan := make(chan *CollatzPair, numWorkers)
	for i := 0; i < numWorkers; i++ {
		go collatzLenProducer(numGen, lenChan)
	}
	return maxOfChan(end-start, lenChan)
}

/** Finds maximum Collatz Length between two numbers in an iterative way
 * @param start the beginning of the sequence to check
 * @param end the end of the sequence non-inclusive
 * @param numWorkers number of worker threads to spawn
 * @return the maximum Collatz number and its length
 */
func maxCollatzIterative(start int, end int) *CollatzPair {
	maxN, maxLen := 1, 1
	for i := start; i < end; i++ {
		length := collatzLen(i)
		if length > maxLen {
			maxN, maxLen = i, length
		}
	}
	return &CollatzPair{maxN, maxLen}
}

func main() {
	numWorkers := flag.Int("workers", 5, "number of collatz workers to spawn")
	flag.Parse()
	//collatzLen(13)
	//timeit("iterative", func() {
	//	maxPair := maxCollatzIterative(1, 10000001)
	//	fmt.Printf("Longest sequence starts at %d, length %d\n",
	//		maxPair.start,
	//		maxPair.length)
	//})
	//timeit("parallel", func() {
	//	maxPair := maxCollatzParallel(1, 10000001, *numWorkers)
	//	fmt.Printf("Longest sequence starts at %d, length %d\n",
	//		maxPair.start,
	//		maxPair.length)
	//})
	maxPair := maxCollatzParallel(1, 10000001, *numWorkers)
	fmt.Printf("Longest sequence starts at %d, length %d\n",
		maxPair.start,
		maxPair.length)
}
