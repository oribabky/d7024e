package d7024e

import (
	"testing"
	"log"
)


const ClientAddress string  = "localhost:7999"
const ClientAddress1 string = "localhost:7998"
const ClientAddress2 string = "localhost:7997"
const ClientAddress3 string = "localhost:7996"
const ClientAddress4 string = "localhost:7995"
const ClientAddress5 string = "localhost:7994"
const ClientAddress6 string = "localhost:7993"
const ClientAddress7 string = "localhost:7992"
const ClientAddress8 string = "localhost:7991"

const ServerAddress1 string = "localhost:8000"
const ServerAddress2 string = "localhost:8001"
const ServerAddress3 string = "localhost:8002"

func TestRPCs(t *testing.T) {
	
	node1 := NewNode("", ClientAddress)
	go node1.NodeUp()



	//test find_node
	log.Println("FIND_NODE")
	go node1.network.SendPingMessage(ServerAddress1)
	node1.network.SendPingMessage(ServerAddress2)
	/*ping3 := go node1.network.SendPingMessage(ServerAddress3)
	ping4 := node1.network.SendPingMessage(ServerAddress1)

	//test ping
	log.Println("PING")
	log.Println(ping1)
	log.Println(ping2)
	log.Println(ping3)
	log.Println(ping4) */

	//node1.network.SendPingMessage(ServerAddress1)

	node1.network.CloseConnection();
	
} 
/*
func TestFindNode(t *testing.T) {

	node1 := NewNode("", ClientAddress)
	go node1.NodeUp()
	go node1.network.SendFindNodeMessage(ServerAddress)
}*/