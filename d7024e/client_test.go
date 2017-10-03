package d7024e

import (
	"testing"
	"log"
	"time"
	//"fmt"
	//"math/rand"
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


/* Test case 2001: The system should be able to send various RPCs.*/
func Test_2001(t *testing.T) {
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
	time.Sleep(time.Millisecond * 500)
	log.Println("\nFIND_NODE")

	target := NewContact(NewRandomKademliaID(), "localhost:8000")

	//add 20 random contacts to each node
	for i := 0; i < 20; i++ {
		server1.Rt.AddContact(NewContact(NewRandomKademliaID(), "localhost:8000"))
		server2.Rt.AddContact(NewContact(NewRandomKademliaID(), "localhost:8000"))
		server3.Rt.AddContact(NewContact(NewRandomKademliaID(), "localhost:8000"))
	}
	
	//we also need to expect that node1 will be in their routing tables because we send requests to them from node1
	server1.Rt.AddContact(*node1.Me)
	server2.Rt.AddContact(*node1.Me)
	server3.Rt.AddContact(*node1.Me)
	node1.Rt.AddContact(*server1.Me)
	node1.Rt.AddContact(*server2.Me)
	node1.Rt.AddContact(*server3.Me) 

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
	        case <-time.After(time.Millisecond * 1000):
		    	log.Println("Channel empty.")
		    	break;

	    	case c := <-node1.network.ReturnedContacts:
	    		kClosestTotalActual = append(kClosestTotalActual, *c)
	    		continue;
	    	}
	    break;
	}
	
	time.Sleep(time.Millisecond * 500)

	//check that the size is the same of both slices
	if len(kClosestTotalActual) != len(kClosestTotalExpected) {
		log.Println(len(kClosestTotalActual))
		log.Println(len(kClosestTotalExpected))
		t.Error("error in testing RPCs. ")
	}

	//check that the contents are the same. That is, check that the returned contacts match the kClosestTotal:
	for i := range kClosestTotalExpected {
		currentContact := &kClosestTotalExpected[i]

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

/* Test case 2002: The system should be able to send out kademlia procedures. */

func Test_2002(t *testing.T) {
	time.Sleep(time.Millisecond * 500)
	log.Println("\nTEST Kademlia procedures..")

	/*node1ID :=   "FFFFFFFF00000000000000000000000000000000";
	server1ID := "1111111200000000000000000000000000000000";
	server2ID := "1111111300000000000000000000000000000000";
	server3ID := "1111111400000000000000000000000000000000";
	server4ID := "1111111500000000000000000000000000000000";
	server5ID := "1111111600000000000000000000000000000000";
	server6ID := "1111111700000000000000000000000000000000";
	server7ID := "1111111800000000000000000000000000000000";
	server8ID := "1111111900000000000000000000000000000000";
	server9ID := "1111111A00000000000000000000000000000000";*/

	//random node ID's
	node1 := NewNode("", ClientAddress)
	server1 := NewNode("", ServerAddress1)
	server2 := NewNode("", ServerAddress2)
	server3 := NewNode("", ServerAddress3)
	server4 := NewNode("", ServerAddress4)
	/*server5 := NewNode("", ServerAddress5)
	server6 := NewNode("", ServerAddress6)
	server7 := NewNode("", ServerAddress7)
	server8 := NewNode("", ServerAddress8)
	server9 := NewNode("", ServerAddress9) */

	/*node1 := NewNode(node1ID, ClientAddress)
	server1 := NewNode(server1ID, ServerAddress1)
	server2 := NewNode(server2ID, ServerAddress2)
	server3 := NewNode(server3ID, ServerAddress3)
	server4 := NewNode(server4ID, ServerAddress4)
	 server5 := NewNode(server5ID, ServerAddress5)
	server6 := NewNode(server6ID, ServerAddress6)
	server7 := NewNode(server7ID, ServerAddress7)
	server8 := NewNode(server8ID, ServerAddress8)
	server9 := NewNode(server9ID, ServerAddress9) */

	go node1.NodeUp()
	go server1.NodeUp()
	go server2.NodeUp()
	go server3.NodeUp()
	go server4.NodeUp()
	/*go server5.NodeUp()
	go server6.NodeUp()
	go server7.NodeUp()
	go server8.NodeUp()
	go server9.NodeUp()  */

	//test nodeLookup
	log.Println("\nNODE LOOKUP")

	//onlineNodes := []*Node{node1, server1, server2, server3, server4, server5, server6, server7, server8, server9}
	onlineNodes := []*Node{node1, server1, server2, server3, server4}
	nrOnlineNodes := len(onlineNodes)
	log.Println(nrOnlineNodes)


	//add every node to node1's routing table, and every node should know of node1:
	for i := range onlineNodes {
		if onlineNodes[i].Me.ID.String() == node1.Me.ID.String() {
			continue;
		}
		node1.Rt.AddContact(*onlineNodes[i].Me)
		onlineNodes[i].Rt.AddContact(*node1.Me)
	} 

	k := K;
	a := Alpha;
	kClosest0 := node1.Kademlia.LookupContact(node1.Me, k, a)
	kClosest1 := server1.Kademlia.LookupContact(server1.Me, k, a)
	kClosest2 := server2.Kademlia.LookupContact(server2.Me, k, a)
	kClosest3 := server3.Kademlia.LookupContact(server3.Me, k, a)
	kClosest4 := server4.Kademlia.LookupContact(server4.Me, k, a)
	/*kClosest5 := server5.Kademlia.LookupContact(server5.Me, k, a)
	kClosest6 := server6.Kademlia.LookupContact(server6.Me, k, a)
	kClosest7 := server7.Kademlia.LookupContact(server7.Me, k, a)
	kClosest8 := server8.Kademlia.LookupContact(server8.Me, k, a)
	kClosest9 := server9.Kademlia.LookupContact(server9.Me, k, a) */

	time.Sleep(time.Millisecond * 1000)
	//kClosestAll := [][]Contact{kClosest0, kClosest1, kClosest2, kClosest3, kClosest4, kClosest5, kClosest6, kClosest7, kClosest8, kClosest9}
	kClosestAll := [][]Contact{kClosest0, kClosest1, kClosest2, kClosest3, kClosest4}

	for i := range onlineNodes {
		log.Println("\nI am node " + onlineNodes[i].Me.Address + " and these are my KClosest:")

		kClosestMe := onlineNodes[i].Rt.FindClosestContacts(onlineNodes[i].Me.ID, k)

		for o := range kClosestAll[i] {
			if kClosestAll[i][o].ID.String() != kClosestMe[o].ID.String() {
				log.Println("Actual: " + kClosestAll[i][o].Address)
				log.Println("Expected: " + kClosestMe[o].Address)
				t.Error("error in test case 2002.")
			}
		}

		onlineNodes[i].Rt.PrintRoutingTable()
	} 



	//from any node we should now be able to find the closest nodes to a given target
	//try to find the closest nodes to server7 from server2
	
	kClosestActual4 := server1.Kademlia.LookupContact(server4.Me, k, a)
	kClosest4 = server4.Rt.FindClosestContacts(server4.Me.ID, k)
	//PrintContactList(kClosest4)
	time.Sleep(time.Millisecond * 1000)
	//log.Println("\nactual closest for node: " + server4.Me.Address)
	for i := range kClosestActual4 {
		if kClosestActual4[i].ID.String() != kClosest4[i].ID.String() {
			log.Println("Target: " + server4.Me.Address + "/" + server4.Me.ID.String())
			log.Println("Actual: " + kClosestActual4[i].Address + "/" + kClosestActual4[i].ID.String())
			log.Println("Expected: " + kClosest4[i].Address + "/" + kClosest4[i].ID.String())
			t.Error("error in test case 2002.")
		}
	}
	



/*
	//connect our nodes to one another through the routing tables. Each node will know of at least 1 other online node.
	for i := range onlineNodes {
		nrActiveNodesToAdd := 1;
		nrActiveNodesToAdd = nrActiveNodesToAdd + rand.Intn(4) 

		//add a certain amount of nodes
		for k := 0; k < nrActiveNodesToAdd; k++ {

			//find an ok random index of the online nodes 
			randomIndexOk := false;
			randomIndex := rand.Intn((nrOnlineNodes - 1))
			for randomIndexOk == false {

				if randomIndex != i {
					randomIndexOk = true
					break;
				}
				randomIndex = rand.Intn((nrOnlineNodes - 1))
			}
			chosenContact := onlineNodes[randomIndex].Me

			//add the node as a contact in the routing table.
			onlineNodes[i].Rt.AddContact(*chosenContact)
		}


	} */

/*
	//our target contact that we want to find k closest contacts to
	target := NewContact(NewKademliaID("1111111500000000000000000000000000000000"), "localhost:8000")

	//our actual kClosest contacts to the target:
	k := 3;
	closest1 := NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:1337")
	closest2 := NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:1338")
	closest3 := NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:1339")
	kClosestExpected := []Contact{closest1, closest2, closest3}
	nrOnlineNodes = len(onlineNodes)
	log.Println(rand.Intn(nrOnlineNodes))
	//add the closest randomly in our online nodes' routing tables.
	for i := range kClosestExpected {
		contact := kClosestExpected[i]
		randomIndex := rand.Intn(nrOnlineNodes - 1)
		log.Println(randomIndex)
		onlineNodes[randomIndex].Rt.AddContact(contact)
	}

	//TEMP print all nodes k-closest.
	for i := range onlineNodes {
		currentNode := onlineNodes[i]
		kClosest := currentNode.Rt.FindClosestContacts(target.ID, 100)
		log.Println("\nHi I am node " + currentNode.Me.Address + ". This is my routing table:")
		for k := range kClosest {
			log.Println(kClosest[k].Address)
		}
	}

	//finally, now we should be able to find the kClosest contacts
	closestToTarget := node1.Kademlia.LookupContact(&target, k)

	
	for i := range closestToTarget {
		log.Println("EXPECTED: " + kClosestExpected[i].ID.String())
		log.Println("ACTUAL " + closestToTarget[i].ID.String())
		// if closestToTarget[i] != kClosestExpected[i] {
		//	t.Error("Error testing nodeLookup..")
		//} 
	}  
	*/
	time.Sleep(time.Millisecond * 500)



	node1.network.CloseConnection();
	server1.network.CloseConnection();
	server2.network.CloseConnection();
	server3.network.CloseConnection();
	server4.network.CloseConnection();
	/* server5.network.CloseConnection();
	server6.network.CloseConnection();
	server7.network.CloseConnection();
	server8.network.CloseConnection();
	server9.network.CloseConnection(); */

	time.Sleep(time.Millisecond * 500)

} 
/*
func TestFindNode(t *testing.T) {

	node1 := NewNode("", ClientAddress)
	go node1.NodeUp()
	go node1.network.SendFindNodeMessage(ServerAddress)
}*/