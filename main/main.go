package main

import (
	d "../d7024e"
)

func main () {
	
	node1 := d.NewNode("", ":4000")
	for i := 0; i < 10; i++ {
		node1.Rt.AddContact(d.NewContact(d.NewRandomKademliaID(), "localhost:8002"))
	}
	node1.NodeUp()

}