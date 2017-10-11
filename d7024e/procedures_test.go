package d7024e

import (
	"testing"
	"log"
	"time"
	//"fmt"
	"math/rand"
	"strconv"
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
//const OfflineServer string = "localhost:9000"


/* Test case 2002: The sytem should be able to locate k-closest nodes to a given target */
func Test_2002(t *testing.T) {
	time.Sleep(time.Millisecond * 500)
	log.Println("\nTEST Kademlia procedures..")

	/*node1ID :=   "FFFFFFFF00000000000000000000000000000000";
	server1ID := "1111111200000000000000000000000000000000";
	server2ID := "1111111300000000000000000000000000000000";
	server3ID := "1111111400000000000000000000000000000000";
	server4ID := "1111111500000000000000000000000000000000";*/


	//random node ID's
	node1 := NewNode("", ClientAddress)
	server1 := NewNode("", ServerAddress1)
	server2 := NewNode("", ServerAddress2)
	server3 := NewNode("", ServerAddress3)
	server4 := NewNode("", ServerAddress4)

	/*node1 := NewNode(node1ID, ClientAddress)
	server1 := NewNode(server1ID, ServerAddress1)
	server2 := NewNode(server2ID, ServerAddress2)
	server3 := NewNode(server3ID, ServerAddress3)
	server4 := NewNode(server4ID, ServerAddress4)*/


	go node1.NodeUp()
	go server1.NodeUp()
	go server2.NodeUp()
	go server3.NodeUp()
	go server4.NodeUp()

	//test nodeLookup
	log.Println("\nNODE LOOKUP")

	onlineNodes := []*Node{node1, server1, server2, server3, server4}
	nrOnlineNodes := len(onlineNodes)

	//add every node to node1's routing table, and every node should know of node1:
	for i := range onlineNodes {
		if onlineNodes[i].Me.ID.String() == node1.Me.ID.String() {
			continue;
		}
		node1.Rt.AddContact(*onlineNodes[i].Me)
		onlineNodes[i].Rt.AddContact(*node1.Me)
	} 


	kClosest0 := node1.Kademlia.LookupContact(node1.Me.ID)
	kClosest1 := server1.Kademlia.LookupContact(server1.Me.ID)
	kClosest2 := server2.Kademlia.LookupContact(server2.Me.ID)
	kClosest3 := server3.Kademlia.LookupContact(server3.Me.ID)
	kClosest4 := server4.Kademlia.LookupContact(server4.Me.ID)

	kClosestAll := [][]Contact{kClosest0, kClosest1, kClosest2, kClosest3, kClosest4}

	for i := range onlineNodes {
		log.Println("\nI am node " + onlineNodes[i].Me.Address + " and these are my KClosest:")

		kClosestMe := onlineNodes[i].Rt.FindClosestContacts(onlineNodes[i].Me.ID, K)

		for o := range kClosestAll[i] {
			if kClosestAll[i][o].ID.String() != kClosestMe[o].ID.String() {
				log.Println("Actual: " + kClosestAll[i][o].Address)
				log.Println("Expected: " + kClosestMe[o].Address)
				t.Fatal("error in test case 2002.")
			}
		}

		onlineNodes[i].Rt.PrintRoutingTable()
	} 



	//from any node we should now be able to find the closest nodes to a given target

	

	//chose at random a source node and a target from which we will try to find the kClosest:
	indexLimit := nrOnlineNodes - 1;

	for i := 0; i < 5; i++ {
		randSourceIndex := rand.Intn(indexLimit)
		randTargetIndex := rand.Intn(indexLimit)
		chosenSourceNode := onlineNodes[randSourceIndex]
		chosenTargetNode := onlineNodes[randTargetIndex]

		kClosestActual := chosenSourceNode.Kademlia.LookupContact(chosenTargetNode.Me.ID)
		kClosestExpected := chosenTargetNode.Kademlia.LookupContact(chosenTargetNode.Me.ID)
		for j := range kClosestActual {
			if kClosestActual[j].ID.String() != kClosestExpected[j].ID.String() {
				log.Println("Target: " + chosenTargetNode.Me.Address + "/" + chosenTargetNode.Me.ID.String())
				log.Println("Actual: " + kClosestActual[j].Address + "/" + kClosestActual[j].ID.String())
				log.Println("Expected: " + kClosestExpected[j].Address + "/" + kClosestExpected[j].ID.String())
				t.Fatal("error in test case 2002.")
			}
		}
	}  



	time.Sleep(time.Millisecond * 500)

	node1.Network.CloseConnection();
	server1.Network.CloseConnection();
	server2.Network.CloseConnection();
	server3.Network.CloseConnection();
	server4.Network.CloseConnection();

	time.Sleep(time.Millisecond * 500)

} 
/*
func TestFindNode(t *testing.T) {

	node1 := NewNode("", ClientAddress)
	go node1.NodeUp()
	go node1.Network.SendFindNodeMessage(ServerAddress)
}*/

/* Test case 2003: The sytem should be store a file in the Network at the k-closest contacts to the file hash. */
func Test_2003(t *testing.T) {
time.Sleep(time.Millisecond * 500)
	//test STORE
	log.Println("\nSTORE")

	//random node ID's
	node1 := NewNode("", ClientAddress)
	server1 := NewNode("", ServerAddress1)
	server2 := NewNode("", ServerAddress2)
	server3 := NewNode("", ServerAddress3)
	server4 := NewNode("", ServerAddress4)


	go node1.NodeUp()
	go server1.NodeUp()
	go server2.NodeUp()
	go server3.NodeUp()
	go server4.NodeUp()

	

	//Fill up the routing tables
	onlineNodes := []*Node{node1, server1, server2, server3, server4}

	//add every node to node1's routing table, and every node should know of node1:
	for i := range onlineNodes {
		if onlineNodes[i].Me.ID.String() == node1.Me.ID.String() {									/* TESTA SÅ ATT NODE1 LÄGGER TILL SIG SJÄLV I RT */
			continue;
		}
		node1.Rt.AddContact(*onlineNodes[i].Me)
		onlineNodes[i].Rt.AddContact(*node1.Me)
	} 

	node1.Kademlia.LookupContact(node1.Me.ID)
	server1.Kademlia.LookupContact(server1.Me.ID)
	server2.Kademlia.LookupContact(server2.Me.ID)
	server3.Kademlia.LookupContact(server3.Me.ID)
	server4.Kademlia.LookupContact(server4.Me.ID)


			//store the file in the system
	fileContents := []byte("asdasdasdasdasd")
	fileKey := node1.Kademlia.Store(fileContents) 

	//kClosestExpected := dummyNode.Kademlia.LookupContact(dummyNode.Me.ID)
	kClosestExpected := node1.Kademlia.LookupContact(fileKey)
	PrintContactList(kClosestExpected)



	time.Sleep(time.Millisecond * 500)

	nrStoredLocationsExpected := K
	if nrStoredLocationsExpected > len(onlineNodes) {
		nrStoredLocationsExpected = len(onlineNodes)
	}

	nrStoredLocationsActual := 0
	actualString := "";

	for i := range onlineNodes {
		if onlineNodes[i].Network.FileExists(fileKey) == true {
			nrStoredLocationsActual ++;
			actualString += onlineNodes[i].Me.Address + "/"
		}		
	}
	
	time.Sleep(time.Millisecond * 500)
	if nrStoredLocationsActual != nrStoredLocationsExpected {
		log.Println(actualString)
		log.Println("Actual: " + strconv.Itoa(nrStoredLocationsActual))
		log.Println("Expected: " + strconv.Itoa(nrStoredLocationsExpected))
		t.Fatal("Error in test case 2003.")
	}


	time.Sleep(time.Millisecond * 500)
	//close the connections
	node1.Network.CloseConnection();
	server1.Network.CloseConnection();
	server2.Network.CloseConnection();
	server3.Network.CloseConnection();
	server4.Network.CloseConnection();
	//dummyNode.Network.CloseConnection();
	time.Sleep(time.Millisecond * 500)
}
/*
/* Test case 2004: The system should be able to locate a file in the Network */
func Test_2004(t *testing.T) {
time.Sleep(time.Millisecond * 500)
	//test STORE
	log.Println("\nCAT")

	//random node ID's
	node1 := NewNode("", ClientAddress)
	server1 := NewNode("", ServerAddress1)
	server2 := NewNode("", ServerAddress2)
	server3 := NewNode("", ServerAddress3)
	server4 := NewNode("", ServerAddress4)

	
	go node1.NodeUp()
	go server1.NodeUp()
	go server2.NodeUp()
	go server3.NodeUp()
	go server4.NodeUp()

	//Fill up the routing tables
	onlineNodes := []*Node{node1, server1, server2, server3, server4}

	//add every node to node1's routing table, and every node should know of node1:
	for i := range onlineNodes {
		if onlineNodes[i].Me.ID.String() == node1.Me.ID.String() {
			continue;
		}
		node1.Rt.AddContact(*onlineNodes[i].Me)
		onlineNodes[i].Rt.AddContact(*node1.Me)
	} 

	node1.Kademlia.LookupContact(node1.Me.ID)
	server1.Kademlia.LookupContact(server1.Me.ID)
	server2.Kademlia.LookupContact(server2.Me.ID)
	server3.Kademlia.LookupContact(server3.Me.ID)
	server4.Kademlia.LookupContact(server4.Me.ID)

	//store the file in the system
	fileContents := []byte("asdasdasdasdasd")
	fileKey := node1.Kademlia.Store(fileContents) 
	log.Println(fileKey.String())

	time.Sleep(time.Millisecond * 500)
	//we should now be able to find the contents of this file from any node
	nrNoFounds := 0
	for i := range onlineNodes {
		actualFileContents := onlineNodes[i].Kademlia.LookupData(fileKey)
		if actualFileContents == nil {
			nrNoFounds++;
		}
		/*if string(actualFileContents) != string(fileContents) {
			t.Fatal("Error in test case 2004.")
		} */
	}
	if nrNoFounds > 0 {
		t.Fatal("Error in test case 2004.")
	}
	log.Println(nrNoFounds)

	time.Sleep(time.Millisecond * 500)
	//close the connections
	node1.Network.CloseConnection();
	server1.Network.CloseConnection();
	server2.Network.CloseConnection();
	server3.Network.CloseConnection();
	server4.Network.CloseConnection();
	time.Sleep(time.Millisecond * 500)
}