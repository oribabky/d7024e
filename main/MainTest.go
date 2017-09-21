package main 

import d "../d7024e"
import "fmt"
import "time"

func main() {
	
	srcNode := d.NewContact(d.NewRandomKademliaID(),"127.0.0.1:8001")
	serverNode := d.NewContact(d.NewRandomKademliaID(),"127.0.0.2:8002")
	differentPort := d.NewContact(d.NewRandomKademliaID(),"127.0.0.1:8003")
	fmt.Println(srcNode.ID)

	// rt := d.NewRoutingTable(srcNode)
	net := d.NewNetwork(&srcNode)
	net2 := d.NewNetwork(&serverNode)
	net3 := d.NewNetwork(&differentPort)
	
	go net2.Listen()
	go net3.Listen()	

	time.Sleep(time.Second * 2)	
	net.SendPingMessage(&serverNode)
	time.Sleep(time.Second * 2)
}