package main

import (
	"bufio"
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
}

type kbngraph map[string]*node

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
	g[name] = &node{nodeType, name, make(map[string]*node), nil}
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
	}
}
