package main

import (
	d "../d7024e"
	"strconv"
	"log"
	"math/rand"
	"fmt"
	"bufio"
	"os"
)

/* This will be a simulation of 100 nodes */
func main () {
	nrNodes := 1;

	//create 100 nodes
	port := 8000;
	nodes := make([]*d.Node, 0)
	for i := 0; i < nrNodes; i++ {
		newNode := d.NewNode("", "localhost:" + strconv.Itoa(port))
		go newNode.NodeUp()

		//the node needs to have at least one other online node in its routing table to connect to the network.
		//unless it's the first node joining the network
		if len(nodes) == 1 {
			newNode.Rt.AddContact(*nodes[0].Me)
		} else if len(nodes) > 1 {
			indexLimit := len(nodes) - 1
			randIndex := rand.Intn(indexLimit)
			randContact := nodes[randIndex].Me
			newNode.Rt.AddContact(*randContact)
		}

		nodes = append(nodes, newNode)
		port ++;

		//the node needs to add itself to the network by performing a lookupContact on it self.
		nodes[i].Kademlia.LookupContact(nodes[i].Me.ID)
	}
	log.Println("Nodes are up")

	for i := range nodes {
		nodes[i].Rt.PrintRoutingTable()
	}
	/*scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Choose procedure: ")
	for scanner.Scan() {

    	fmt.Println(scanner.Text())
	}*/
	
	for {
		reader := bufio.NewReader(os.Stdin)
	    fmt.Print("Choose procedure: ")
	    text, _ := reader.ReadString('\n')
	    switch text {
		    case "store":
		    case "cat":
		    case "pin":
		    case "unpin":
		    default:
	    }


	    
	    fmt.Println(text[0])

	    fmt.Println("Enter text: ")
	    text2 := ""
	    fmt.Scanln(text2)
	    fmt.Println(text2)

	    ln := ""
	    fmt.Sscanln("%v", ln)
	    fmt.Println(ln)	
		} 

}

