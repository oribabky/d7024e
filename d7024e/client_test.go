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
const ServerAddress4 string = "localhost:8003"
const OfflineServer string = "localhost:9000"

func TestRPCs(t *testing.T) {
	
	node1 := NewNode("", ClientAddress)
	go node1.NodeUp()


	//test ping
	log.Println("PING")

	node1.network.SendPingMessage(ServerAddress1)
	node1.network.SendPingMessage(ServerAddress2)
    node1.network.SendPingMessage(ServerAddress3)
	node1.network.SendPingMessage(ServerAddress4)
	node1.network.SendPingMessage(OfflineServer)



	//test find_node, should be able to be sent asynchronously.
	log.Println("FIND_NODE")
	go node1.network.SendFindNodeMessage(ServerAddress1, node1.Me.ID.String())
	go node1.network.SendFindNodeMessage(ServerAddress1, node1.Me.ID.String())
	go node1.network.SendFindNodeMessage(ServerAddress1, node1.Me.ID.String())

	for {
		c := <- node1.network.ReturnedContacts
		log.Println(c.ID.String())
	}
	//node1.network.SendPingMessage(ServerAddress1)

	
	

	for {

	} 
	node1.network.CloseConnection();
} 
/*
func TestFindNode(t *testing.T) {

	node1 := NewNode("", ClientAddress)
	go node1.NodeUp()
	go node1.network.SendFindNodeMessage(ServerAddress)
}*/