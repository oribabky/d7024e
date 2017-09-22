package main

import (
	d "../d7024e"
	"log"
	)

func main() {
	srcNode := "FFFFFFFF00000000000000000000000000000000";
	server := d.NewContact(d.NewKademliaID(srcNode), ":4000");
	network := d.NewNetwork(&server)

	log.Println("Hello I am on address: " + server.Address)
	go network.Listen()
	network.RequestHandler()
}