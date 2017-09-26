package d7024e

import (
	"testing"
	"log"
)

const ClientAddress string  = ":4001"
const ServerAddress string = ":4000"

func TestRPCs(t *testing.T) {
	
	node1 := NewNode("", ClientAddress)
	go node1.NodeUp()

	//test ping
	log.Println("---PING---")
	node1.network.SendPingMessage(ServerAddress)

	log.Println("---FIND_NODE---")
	//test find_node
	node1.network.SendFindNodeMessage(ServerAddress, node1.Me.ID.String())
	for {

	}
}
/*
func TestFindNode(t *testing.T) {

	node1 := NewNode("", ClientAddress)
	go node1.NodeUp()
	go node1.network.SendFindNodeMessage(ServerAddress)
}*/