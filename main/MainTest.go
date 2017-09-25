package main 

import d "../d7024e"
import "fmt"
import "time"

func main() {
	srcNode := d.NewContact(d.NewRandomKademliaID(),"127.0.0.1:8001")
	Kademlia2 := d.NewKademlia(&srcNode)
	net := d.NewNetwork(&srcNode)

	d.SendFindContactMessage(&srcNode, &Kademlia2)

}