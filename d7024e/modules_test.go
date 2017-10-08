package d7024e

import (
	"testing"
	"fmt"
	"log"
	"math/rand"
	"time"
)


/* Test case 1001: FindClosestContacts should return k closest nodes ordered in distance to source node. */
func TestRoutingTable_1001(t *testing.T) {
	time.Sleep(time.Millisecond * 500)
	log.Println("testing")
	node0 := NewNode("", ClientAddress)
	srcNode := "FFFFFFFF00000000000000000000000000000000";
	node1 := "1111111100000000000000000000000000000000";
	node2 := "1111111200000000000000000000000000000000";
	node3 := "1111111300000000000000000000000000000000";
	node4 := "1111111400000000000000000000000000000000";
	node5 := "2111111400000000000000000000000000000000";

	//expected k closest in the order they should they appear.
	expected := [5]string{node5, node4, node3, node2, node1}

	rt := node0.Rt;
	rt.AddContact(NewContact(NewKademliaID(node1), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID(node2), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID(node3), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID(node4), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID(node5), "localhost:8002"))

	k := 5;
	kClosest := rt.FindClosestContacts(NewKademliaID(srcNode), k)

	for i := range kClosest {
		if kClosest[i].ID.String() != expected[i] {
			t.Fatal("error in test case 1001")
		}
	}
	time.Sleep(time.Millisecond * 500)
	node0.Network.CloseConnection();
	time.Sleep(time.Millisecond * 500)
}

/* Test case 1002: FindClosestContact should return k items even if k is higher than the amount of items available. */
func TestRoutingTable_1002(t *testing.T) {

	node1 := "1111111100000000000000000000000000000000";
	node2 := "1111111200000000000000000000000000000000";
	node3 := "1111111300000000000000000000000000000000";

	node0 := NewNode("", ClientAddress)
	rt := node0.Rt;

	rt.AddContact(NewContact(NewKademliaID(node2), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID(node3), "localhost:8002"))

	k := 1000;
	kClosest := rt.FindClosestContacts(NewKademliaID(node1), k)
    fmt.Println(kClosest[0].Address)
	expectedSize := 3;

	if len(kClosest) != expectedSize {
			t.Fatal("error in test case 1002")
		}
	time.Sleep(time.Millisecond * 500)	
	node0.Network.CloseConnection();
	time.Sleep(time.Millisecond * 500)

}

/* Test case 1003: AddContact should add a contact to the routing table only if that ID is not already taken, regardless of port used.*/
func TestRoutingTable_1003(t *testing.T) {
	node1 := "1111111100000000000000000000000000000000";
	node2 := "1111111200000000000000000000000000000000";

	node0 := NewNode("", ClientAddress)
	rt := node0.Rt;

	rt.AddContact(NewContact(NewKademliaID(node1), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID(node2), "localhost:8002"))

	//try adding same ID as node2 again with same port
	expectedNrNodes := 3;
	rt.AddContact(NewContact(NewKademliaID(node2), "localhost:8002"))
	nodes := rt.FindClosestContacts(NewKademliaID(node1), 100)

	if len(nodes) != expectedNrNodes {
			t.Fatal("error in test case 1003")
		}

	//Adding a node with the same ID as another in the routing table but a new port should not change the routing table.
	newPort := "localhost:8303"
	rt.AddContact(NewContact(NewKademliaID(node2), newPort))
	nodes = rt.FindClosestContacts(NewKademliaID(node1), 100)

	if len(nodes) != expectedNrNodes {
		t.Fatal("error in test case 1003")
	}

	if nodes[1].Address == newPort {
		t.Fatal("error in test case 1003")
	}
	
	node0.Network.CloseConnection();
}

/* Test case 1004: When adding a contact to the routing table bucket that is full, the system should ping the least
	recently seen node in the bucket. If it responds, then nothing should be done. Otherwise if the node doesnt respond
	it should be evicted from the bucket and the new contact should be added. */
func testRoutingTable_1004(t *testing.T) {	

	node1 := NewNode ("1111111100000000000000000000000000000000", ClientAddress);
	node2 := NewNode ("1111111200000000000000000000000000000000", ClientAddress1);
	node3 := NewNode ("1111111300000000000000000000000000000000", ClientAddress2);
	node4 := NewNode ("1111111400000000000000000000000000000000", ClientAddress3);
	node5 := NewNode ("1111111500000000000000000000000000000000", ClientAddress4);
	node6 := NewNode ("1111111600000000000000000000000000000000", ClientAddress5);
	node7 := NewNode ("1111111700000000000000000000000000000000", ClientAddress6);
	node8 := NewNode ("1111111800000000000000000000000000000000", ClientAddress7);
	node9 := NewNode ("1111111900000000000000000000000000000000", ClientAddress8);

	go node1.NodeUp()
	go node2.NodeUp()
	go node3.NodeUp()
	go node4.NodeUp()
 	go node5.NodeUp()
	go node6.NodeUp()
	go node7.NodeUp()
	go node8.NodeUp()
	go node9.NodeUp() 

	node1.Rt.AddContact(*node2.Me)
	node1.Rt.AddContact(*node3.Me)
 	node1.Rt.AddContact(*node4.Me)
	node1.Rt.AddContact(*node5.Me)
	node1.Rt.AddContact(*node6.Me) 
	node1.Rt.AddContact(*node7.Me)
	node1.Rt.AddContact(*node8.Me) 
	node1.Rt.AddContact(*node9.Me)


	/*contacts := node1.Rt.FindClosestContacts(node2.Me.ID, 20)
	for i := range contacts {
		log.Println(contacts[i].ID.String())
		log.Println("bucket index: " + strconv.Itoa(node1.Rt.getBucketIndex(contacts[i].ID)))
	}

	chosenBucketIndex := 30;
	for e := node1.Rt.buckets[chosenBucketIndex].list.Front(); e != nil; e = e.Next() {
		log.Println(e.Value.(Contact).ID.String())
	}

	//now if we try to insert node7 to the routing table:
	

	log.Println("\n")
	for e := node1.Rt.buckets[chosenBucketIndex].list.Front(); e != nil; e = e.Next() {
		log.Println(e.Value.(Contact).ID.String())
	} */

	node1.Network.CloseConnection();
	node2.Network.CloseConnection();
	node3.Network.CloseConnection();
	node4.Network.CloseConnection();
	node5.Network.CloseConnection();
	node6.Network.CloseConnection();
	node7.Network.CloseConnection();
	node8.Network.CloseConnection();
	node9.Network.CloseConnection();


}

/* Test case 1005: When calling the system calls "InsertContactSortedDistTarget" it should try to insert
an item into the right place in a list sorted on distance to a certain target. */
func Test_1005(t *testing.T) {

	targetID := "FFFFFFFF00000000000000000000000000000000";
	node1ID :="1111111100000000000000000000000000000000";
	node2ID := "1111111300000000000000000000000000000000";
	node3ID := "1111111500000000000000000000000000000000";
	node4ID := "1111111700000000000000000000000000000000";
	node5ID := "2111111900000000000000000000000000000000";

	target := NewContact(NewKademliaID(targetID), "localhost:8002")
	node1 := NewContact(NewKademliaID(node1ID), "localhost:8002")
	node2 := NewContact(NewKademliaID(node2ID), "localhost:8002")
	node3 := NewContact(NewKademliaID(node3ID), "localhost:8002")
	node4 := NewContact(NewKademliaID(node4ID), "localhost:8002")
	node5 := NewContact(NewKademliaID(node5ID), "localhost:8002")

	contacts := []Contact{node5, node4, node3, node2}


	//adding node1 should place it in the end:
	expected := []Contact{node5, node4, node3, node2, node1}
	contacts = InsertContactSortedDistTarget(&node1, contacts, target.ID)

	for i := range expected {
		if expected[i].ID.String() != contacts[i].ID.String() {
			t.Fatal("error in test case 1005")
		}
	}

	//adding node5 should place it in the beginning:
	expected = []Contact{node5, node4, node3, node2}
	contacts = []Contact{node4, node3, node2}

	contacts = InsertContactSortedDistTarget(&node5, contacts, target.ID)

	for i := range expected {
		if expected[i].ID.String() != contacts[i].ID.String() {
			t.Fatal("error in test case 1005")
		}
	}

	//adding node3 should place it in the middle:
	expected = []Contact{node5, node4, node3, node2, node1}
	contacts = []Contact{node5, node4, node2, node1}

	contacts = InsertContactSortedDistTarget(&node3, contacts, target.ID)

	for i := range expected {
		if expected[i].ID.String() != contacts[i].ID.String() {
			t.Fatal("error in test case 1005")
		}
	}





}

/* Test case 1006: ContainsContact should tell us whether a contact exists in a list of contacts or not. */
func Test_1006(t *testing.T) {

	contacts := make([]Contact, 0)

	nrContacts := 10;
	for i := 0; i < nrContacts; i++ {
		contacts = append(contacts, NewContact(NewRandomKademliaID(),"localhost:8000"))
	}

	//randIndex := 5; //temp	
	randIndex := rand.Intn(nrContacts - 1)

	contact1 := contacts[randIndex]

	if ContainsContact(contacts, &contact1) == false {
		t.Error("error in test case 1006")
	}

	

	contact1 = NewContact(NewKademliaID("e027a259826185ec8aec6b45cd861fca0f22cf6f"),"localhost:8000")

	contacts = append(contacts, contact1)

	if ContainsContact(contacts, &contact1) == false {
		t.Error("error in test case 1006")
	}
}


