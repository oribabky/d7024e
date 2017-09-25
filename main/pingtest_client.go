package main

import (
	d "../d7024e"
	"log"
	)

func main() {
	srcNode := "FFFFFFFF00000000000000000000000000000000";
	me := d.NewContact(d.NewKademliaID(srcNode), ":4001");
	target := d.NewContact(d.NewKademliaID(srcNode), ":4000");

	rt := d.NewRoutingTable(me)
	kademlia := d.NewKademlia(&me, rt)
	network := d.NewNetwork(&me, &kademlia)

	log.Println("Hello I am on address: " + me.Address + " and I will try to ping: " + target.Address)
	go network.Listen()
	//go network.SendKademliaPacket(&target, "pingRequest")
	go network.SendKademliaPacket(&target, "findNodeRequest", "1111111300000000000000000000000000000000")
	network.RequestHandler()

	
}