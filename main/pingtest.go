package main

import d "../d7024e"

func main() {
	srcNode := "FFFFFFFF00000000000000000000000000000000";
	server := d.NewContact(d.NewKademliaID(srcNode), ":4000");
	network := d.NewNetwork(&server)
	network.Listen()
}