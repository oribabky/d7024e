package d7024e

import (
	"testing"
	"log"
	"time"
	//"fmt"
	"math/rand"
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
const ServerAddress5 string = "localhost:8004"
const ServerAddress6 string = "localhost:8005"
const ServerAddress7 string = "localhost:8006"
const ServerAddress8 string = "localhost:8007"
const ServerAddress9 string = "localhost:8008"
const OfflineServer string = "localhost:9000"


/* Test case 2005: The system should be able to send various RPCs.*/
func Test_2005(t *testing.T) {
	log.Println("TEST RPCS..")

	node1 := NewNode("", ClientAddress)
	server1 := NewNode("", ServerAddress1)
	server2 := NewNode("", ServerAddress2)
	server3 := NewNode("", ServerAddress3)
	go node1.NodeUp()
	go server1.NodeUp()
	go server2.NodeUp()
	go server3.NodeUp()

	//test ping
	log.Println("\nPING")

	if node1.network.SendPingMessage(ServerAddress1) == false {
		t.Error("error in testing RPCs.")
	}
	if node1.network.SendPingMessage(ServerAddress2) == false {
		t.Error("error in testing RPCs.")
	}
	if node1.network.SendPingMessage(ServerAddress3) == false {
		t.Error("error in testing RPCs.")
	}
	if node1.network.SendPingMessage(OfflineServer) == true {
		t.Error("error in testing RPCs.")
	}



	//test find_node, should be able to be sent asynchronously.
	time.Sleep(time.Second * 1)
	log.Println("\nFIND_NODE")

	target := NewContact(NewRandomKademliaID(), "localhost:8000")

	//add 20 random contacts to each node
	for i := 0; i < 20; i++ {
		server1.Rt.AddContact(NewContact(NewRandomKademliaID(), "localhost:8000"))
		server2.Rt.AddContact(NewContact(NewRandomKademliaID(), "localhost:8000"))
		server3.Rt.AddContact(NewContact(NewRandomKademliaID(), "localhost:8000"))
	}
	
	//fetch kClosest for each node to the target ID
	kClosestExpected1 := server1.Rt.FindClosestContacts(target.ID, K)
	kClosestExpected2 := server2.Rt.FindClosestContacts(target.ID, K)
	kClosestExpected3 := server3.Rt.FindClosestContacts(target.ID, K)

	kClosestTotalExpected := kClosestExpected1
	kClosestTotalExpected = append(kClosestTotalExpected, kClosestExpected2...)
	kClosestTotalExpected = append(kClosestTotalExpected, kClosestExpected3...)
	kClosestTotalExpected = append(kClosestTotalExpected, kClosestExpected1...)
	kClosestTotalExpected = append(kClosestTotalExpected, kClosestExpected2...)
	kClosestTotalExpected = append(kClosestTotalExpected, kClosestExpected3...)

	//send out the find_node rpcs asynchronously
	go node1.network.SendFindNodeMessage(ServerAddress1, target.ID.String())
	go node1.network.SendFindNodeMessage(ServerAddress2, target.ID.String())
	go node1.network.SendFindNodeMessage(ServerAddress3, target.ID.String())
	go node1.network.SendFindNodeMessage(ServerAddress1, target.ID.String())
	go node1.network.SendFindNodeMessage(ServerAddress2, target.ID.String())
	go node1.network.SendFindNodeMessage(ServerAddress3, target.ID.String())


	kClosestTotalActual := make([]Contact, 0)
	//fetch the returned contacts
	for {
		select {
	        case <-time.After(time.Second * 1):
		    	log.Println("Channel empty.")
		    	break;

	    	case c := <-node1.network.ReturnedContacts:
	    		kClosestTotalActual = append(kClosestTotalActual, *c)
	    		continue;
	    	}
	    break;
	}
	
	time.Sleep(time.Second * 1)

	//check that the size is the same of both slices
	if len(kClosestTotalActual) != len(kClosestTotalExpected) {
		t.Error("error in testing RPCs.")
	}

	//check that the contents are the same. That is, check that the returned contacts match the kClosestTotal:
	for i := range kClosestTotalExpected {
		currentContact := kClosestTotalExpected[i]

		foundMatch := ContainsContact(kClosestTotalActual, currentContact)

		if foundMatch == false {
			t.Error("error in testing RPCs.")
		}
	}

	node1.network.CloseConnection();
	server1.network.CloseConnection();
	server2.network.CloseConnection();
	server3.network.CloseConnection();
} 

/* Test case 2005: The system should be able to send out kademlia procedures. */

func Test_2005(t *testing.T) {
	log.Println("\nTEST Kademlia procedures..")

	node1 := NewNode("", ClientAddress)
	server1 := NewNode("", ServerAddress1)
	server2 := NewNode("", ServerAddress2)
	server3 := NewNode("", ServerAddress3)
	server4 := NewNode("", ServerAddress4)
	server5 := NewNode("", ServerAddress5)
	server6 := NewNode("", ServerAddress6)
	server7 := NewNode("", ServerAddress7)
	server8 := NewNode("", ServerAddress8)
	server9 := NewNode("", ServerAddress9)

	go node1.NodeUp()
	go server1.NodeUp()
	go server2.NodeUp()
	go server3.NodeUp()
	go server4.NodeUp()
	go server5.NodeUp()
	go server6.NodeUp()
	go server7.NodeUp()
	go server8.NodeUp()
	go server9.NodeUp()

	//test nodeLookup
	time.Sleep(time.Second * 1)
	log.Println("\nNODE LOOKUP")

	//add some random contacts to each node's routing table
	for i := 0; i < 10; i++ {
		node1.Rt.AddContact(NewContact(NewRandomKademliaID(), "localhost:8000"))
		server1.Rt.AddContact(NewContact(NewRandomKademliaID(), "localhost:8000"))
		server2.Rt.AddContact(NewContact(NewRandomKademliaID(), "localhost:8000"))
		server3.Rt.AddContact(NewContact(NewRandomKademliaID(), "localhost:8000"))
		server4.Rt.AddContact(NewContact(NewRandomKademliaID(), "localhost:8000"))
		server5.Rt.AddContact(NewContact(NewRandomKademliaID(), "localhost:8000"))
		server6.Rt.AddContact(NewContact(NewRandomKademliaID(), "localhost:8000"))
		server7.Rt.AddContact(NewContact(NewRandomKademliaID(), "localhost:8000"))
		server8.Rt.AddContact(NewContact(NewRandomKademliaID(), "localhost:8000"))
		server9.Rt.AddContact(NewContact(NewRandomKademliaID(), "localhost:8000"))
	}

	onlineNodes := []Node{node1, server1, server2, server3, server4, server5, server6, server7, server8, server9}

	//connect our nodes to one another through the routing tables. Each node will know of at least 1 other online node.
	for i := range onlineNodes {

	}

	//our target contact that we want to find k closest contacts to
	target := NewContact(NewKademliaID("1111111500000000000000000000000000000000"), "localhost:8000")

	//our actual kClosest contacts to the target:
	k := 3;
	closest1 := NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8000")
	closest2 := NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8000")
	closest3 := NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8000")
	kClosestExpected := []Contact{closest1, closest2, closest3}

	//add the closest contacts in 
	closestToTarget := node1.Kademlia.LookupContact()

	node1.network.CloseConnection();
	server1.network.CloseConnection();
	server2.network.CloseConnection();
	server3.network.CloseConnection();
	server4.network.CloseConnection();
	server5.network.CloseConnection();
	server6.network.CloseConnection();
	server7.network.CloseConnection();
	server8.network.CloseConnection();
	server9.network.CloseConnection();

} 
/*
func TestFindNode(t *testing.T) {

	node1 := NewNode("", ClientAddress)
	go node1.NodeUp()
	go node1.network.SendFindNodeMessage(ServerAddress)
}*/