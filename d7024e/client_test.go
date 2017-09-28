package d7024e

import (
	"testing"
	"log"
)

const ClientAddress string  = ":7999"
const ServerAddress1 string = ":8000"
const ServerAddress2 string = ":8001"
const ServerAddress3 string = ":8002"

func TestRPCs(t *testing.T) {
	
	node1 := NewNode("", ClientAddress)
	go node1.NodeUp()



	//test find_node
//	log.Println("FIND_NODE")
	go node1.network.SendPingMessage(ServerAddress1)
	go node1.network.SendPingMessage(ServerAddress2)
	go node1.network.SendPingMessage(ServerAddress3)
	//node1.network.SendPingMessage(ServerAddress1)

	//test ping
	log.Println("PING")
	//node1.network.SendPingMessage(ServerAddress1)
	for {

	}
}
/*
func TestFindNode(t *testing.T) {

	node1 := NewNode("", ClientAddress)
	go node1.NodeUp()
	go node1.network.SendFindNodeMessage(ServerAddress)
}*/