package main

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"log"
	"os"
)

const (
	actorType int = iota
	movieType
)

type node struct {
	nodeType  int
	data      string
	neighbors map[string]*node
	prev      *node // Used for pathing back to Kevin Bacon

	priority int // For use by priority queue
	index    int
}

type kbngraph map[string]*node

// Priority queue implementation adapted from:
// https://golang.org/pkg/container/heap/
type PriorityQueue []*node

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*node)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	if n == 0 {
		return nil
	}
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) update(item *node, priority int) {
	item.priority = priority
	heap.Fix(pq, item.index)
}

/** To string for node for easier debugging
* @return returns a string of the list of neighbors of the node
 */
func (n *node) String() string {
	neighborNames := make([]string, len(n.neighbors))
	for _, neighbor := range n.neighbors {
		neighborNames = append(neighborNames, neighbor.data)
	}
	return fmt.Sprintf("%v", neighborNames)
}

/** Reconstructs path and prints path to Kevin Bacon
* This assumes that Dijkstra's algorithm has already been run on the graph
* @return Kevin Bacon Number of node n
 */
func (n *node) PrintPath() int {
	if n.data == "Kevin Bacon" {
		return 0
	} else if n.prev == nil {
		return -1
	}
	movie := n.prev
	// This always exists because otherwise the path would not exist
	next := movie.prev
	fmt.Printf("%s was in %s with %s\n", n.data, movie.data, next.data)
	return next.PrintPath() + 1
}

/** Dijkstra shortest path implementation, runs on graph and tags all nodes
* @param src the origin node where dijkstra starts from
*/
func (g kbngraph) Dijkstra(src string) {
	unvisited := make(PriorityQueue, len(g))
	i := 0
	for k, v := range g {
		v.priority = len(g)
		if k == src {
			v.priority = 0
		}
		v.index = i
		unvisited[i] = v
		i++
	}
	heap.Init(&unvisited)

	for len(unvisited) > 0 {
		cur := heap.Pop(&unvisited).(*node)
		for _, neighbor := range cur.neighbors {
			if neighbor.priority > cur.priority+1 {
				neighbor.prev = cur
				unvisited.update(neighbor, cur.priority+1)
			}
		}
	}
}

/** Checks if a node exists and if not constructs it and adds it to the graph
* @param nodeType integer identifier of the type of node
* @param name string identifier of the node
* @return pointer to the node
 */
func (g kbngraph) AddNode(nodeType int, name string) *node {
	n, exists := g[name]
	if exists {
		return n
	}
	g[name] = &node{nodeType, name, make(map[string]*node), nil, -1, -1}
	return g[name]
}

/** Takes two nodes and connects them together by adding them as neighbors
* @param id1 identifier of the first node
* @param id2 identifier of the second node
* @return returns an error if the connection could not be made
 */
func connect(n1 *node, n2 *node) error {
	if n1 == nil || n2 == nil {
		return errors.New("cannot connect nil nodes")
	}
	n1.neighbors[n2.data], n2.neighbors[n1.data] = n2, n1
	return nil
}

func main() {
	fin, err := os.Open("cast.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fin.Close()

    // Construction of graph and preprocessing
	g := make(kbngraph)
	scanner := bufio.NewScanner(fin)
	for scanner.Scan() { // Assumes first line is a movie title
		// Construct movie node
		movieName := scanner.Text()
		movie := g.AddNode(movieType, movieName)
		// Takes all of the actors and connects them to the movie
		for scanner.Scan() {
			name := scanner.Text()
			if name == "" {
				break
			}
			actor := g.AddNode(actorType, name)
			connect(movie, actor)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	g.Dijkstra("Kevin Bacon") // Tags graph for all shortest paths to Kevin
	fmt.Println("Loading complete!")

	stdinScanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter actor name: ")
		stdinScanner.Scan()
		actorName := stdinScanner.Text()
		actor, exists := g[actorName]
		if !exists {
			fmt.Println("Unknown actor name")
		} else {
			kbn := actor.PrintPath()
			if kbn == -1 {
				fmt.Println("Infinite KBN")
			} else {
				fmt.Println("Found with KBN of", kbn)
			}
		}
		fmt.Println()
	}
}
