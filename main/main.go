package main

import (
	d "../d7024e"
	"strconv"
	"log"
	"math/rand"
	"fmt"
	"bufio"
	"os"
	"strings"
	"time"
)

/* This will be a simulation of 100 nodes */
func main () {
	nrNodes := 100;

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
	
	for {
		reader := bufio.NewReader(os.Stdin)
	    fmt.Println("\nChoose what to do in the following way separated by blankspace: ")
	    fmt.Println("/Choose node: [0-" + strconv.Itoa(nrNodes - 1) + "]/Choose procedure: store, cat, pin, unpin/If store: file contents, Else: 40 char key")
	    fmt.Print("Command: ")
	    text, _ := reader.ReadString('\n')

	    words := strings.Fields(text)
	    if len(words) != 3 {
	    	fmt.Println("Too few arguments..")
	    	continue;
	    }

	    option1, err := strconv.Atoi(words[0])
	    if err != nil {
	    	fmt.Println("Bad format of node number..")
	    	continue;
	    }
	    if option1 < 0 || option1 > nrNodes - 1 {
	    	fmt.Println("Invalid node chosen..")
	    	continue;
	    } 

	    option2 := words[1]
	    if option2 != "store" && option2 != "cat" && option2 != "pin" && option2 != "unpin" {
	    	fmt.Println(option2)
	    	fmt.Println("Bad procedure choice..")
	    	continue;
	    }

	    option3 := words[2]
	    allowedChars := "0123456789abcdef"
	    switch option2 {
	    case "store":
	    	fileContents := []byte(option3)
	    	ID := nodes[option1].Kademlia.Store(fileContents)
	    	time.Sleep(time.Millisecond * 500)
	    	fmt.Println(ID.String())
	    	continue;

    	default:
    		if len(option3) != 40 {
    			fmt.Println("Must be 40 chars")    		
    			continue;
    		}

    		for i := range option3 {
    			characterFound := false
    			for o := range allowedChars {
    				if option3[i] == allowedChars[o] {
    					characterFound = true;
    					break;
    				}
    			}
    			if characterFound == false {
    				fmt.Println("Non-allowed characters.")
    			}
    		}

	    }


	    kademliaID := d.NewKademliaID(option3)

	    switch option2 {
	    case "cat":
	    	data := nodes[option1].Kademlia.LookupData(kademliaID)
	    	time.Sleep(time.Millisecond * 500)
	    	fmt.Println("File contents: "string(data))

	    case "pin":
	    	nodes[option1].Network.Pin(kademliaID)
	    	time.Sleep(time.Millisecond * 500)
    	
    	case "unpin":
    		nodes[option1].Network.UnPin(kademliaID)
    		time.Sleep(time.Millisecond * 500)
	    }

	


	
	} 

}

