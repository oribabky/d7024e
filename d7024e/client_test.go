package d7024e

import (
	"testing"
	"log"
)

const ClientAddress string  = "localhost:7999"
const ServerAddress1 string = "localhost:8000"
const ServerAddress2 string = "localhost:8001"
const ServerAddress3 string = "localhost:8002"

func TestRPCs(t *testing.T) {
	
	node1 := NewNode("", ClientAddress)
	go node1.NodeUp()



	//test find_node
	log.Println("FIND_NODE")
	go node1.network.SendPingMessage(ServerAddress1)
	go node1.network.SendPingMessage(ServerAddress2)
	/*ping3 := go node1.network.SendPingMessage(ServerAddress3)
	ping4 := node1.network.SendPingMessage(ServerAddress1)

	//test ping
	log.Println("PING")
	log.Println(ping1)
	log.Println(ping2)
	log.Println(ping3)
	log.Println(ping4) */

	//node1.network.SendPingMessage(ServerAddress1)


	
}
/*
func TestFindNode(t *testing.T) {

	node1 := NewNode("", ClientAddress)
	go node1.NodeUp()
	go node1.network.SendFindNodeMessage(ServerAddress)
}*/