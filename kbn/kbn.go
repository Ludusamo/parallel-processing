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
}

type kbngraph map[string]*node

func (n node) String() string {
    neighborNames := make([]string, len(n.neighbors))
    for _, neighbor := range n.neighbors {
        neighborNames = append(neighborNames, neighbor.data)
    }
    return fmt.Sprintf("%v", neighborNames)
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
	g[name] = &node{nodeType, name, make(map[string]*node)}
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
		fmt.Println(g)
		break
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
