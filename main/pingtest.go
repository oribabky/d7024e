package main

import (
	d "../d7024e"
	"log"
	)

func main() {
	srcNode := "FFFFFFFF00000000000000000000000000000000";
	node1 := "1111111100000000000000000000000000000000";
	node2 := "1111111200000000000000000000000000000000";

	server := d.NewContact(d.NewKademliaID(srcNode), ":4000");
	rt := d.NewRoutingTable(server)
	kademlia := d.NewKademlia(&server, rt)
	network := d.NewNetwork(&server, &kademlia)

	rt.AddContact(d.NewContact(d.NewKademliaID(node1), "localhost:8002"))
	rt.AddContact(d.NewContact(d.NewKademliaID(node2), "localhost:8002"))

	log.Println("Hello I am on address: " + server.Address)
	go network.Listen()
	network.RequestHandler()
}