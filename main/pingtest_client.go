package main

import d "../d7024e"

func main() {

	me := d.NewContact(d.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:0");
	target := d.NewContact(d.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:4000");

	network := d.NewNetwork(&me)

	//network.Listen()
	//network.RequestHandler()

	network.SendKademliaPacket(&target, "pingRequest")
}