package main

import d "../d7024e"

func main() {
	srcNode := "FFFFFFFF00000000000000000000000000000000";
	client := d.NewContact(d.NewKademliaID(srcNode), "127.0.0.1:0");
	server := d.NewContact(d.NewKademliaID(srcNode), "127.0.0.1:4000");

	network := d.NewNetwork(&client)

	network.SendPingMessage(&server)
}