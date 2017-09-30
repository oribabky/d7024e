package main

import (
	d "../d7024e"
)

func main () {
	
	node1 := d.NewNode("", "localhost:8000")
	node2 := d.NewNode("", "localhost:8001")
	node3 := d.NewNode("", "localhost:8002")
	node4 := d.NewNode("", "localhost:8003")


	go node1.NodeUp()
	go node2.NodeUp()
	go node3.NodeUp()
	go node4.NodeUp()

	for i := 0; i < 10; i++ {
		node1.Rt.AddContact(d.NewContact(d.NewRandomKademliaID(), "localhost:8002"))
		node2.Rt.AddContact(d.NewContact(d.NewRandomKademliaID(), "localhost:8002"))
		node3.Rt.AddContact(d.NewContact(d.NewRandomKademliaID(), "localhost:8002"))
	}

	for {

	}
	

}