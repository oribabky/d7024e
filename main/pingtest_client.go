package main

import (
	d "../d7024e"
	"log"
	)

func main() {
	srcNode := "FFFFFFFF00000000000000000000000000000000";
	me := d.NewContact(d.NewKademliaID(srcNode), ":4001");
	target := d.NewContact(d.NewKademliaID(srcNode), ":4000");

	network := d.NewNetwork(&me)

	log.Println("Hello I am on address: " + me.Address + " and I will try to ping: " + target.Address)
	go network.Listen()
	go network.SendKademliaPacket(&target, "pingRequest")
	network.RequestHandler()

	
}