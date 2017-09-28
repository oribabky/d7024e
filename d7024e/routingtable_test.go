package d7024e

import (
	"testing"
	"fmt"
	"log"
	//"strconv"
)

/* Test case 1001: FindClosestContacts should return k closest nodes ordered in distance to source node. */
func TestRoutingTable_1001(t *testing.T) {
	
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
			t.Error("error in test case 1001")
		}
	}


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
	expectedSize := 2;

	if len(kClosest) != expectedSize {
			t.Error("error in test case 1002")
		}
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
	expectedNrNodes := 2;
	rt.AddContact(NewContact(NewKademliaID(node2), "localhost:8002"))
	nodes := rt.FindClosestContacts(NewKademliaID(node1), 100)

	if len(nodes) != expectedNrNodes {
			t.Error("error in test case 1003")
		}

	//Adding a node with the same ID as another in the routing table but a new port should not change the routing table.
	expectedNrNodes = 2;
	newPort := "localhost:8303"
	rt.AddContact(NewContact(NewKademliaID(node2), newPort))
	nodes = rt.FindClosestContacts(NewKademliaID(node1), 100)

	if len(nodes) != expectedNrNodes {
		t.Error("error in test case 1003")
	}

	if nodes[1].Address == newPort {
		t.Error("error in test case 1003")
	}
}

/* Test case 1004: When adding a contact to the routing table bucket that is full, the system should ping the least
	recently seen node in the bucket. If it responds, then nothing should be done. Otherwise if the node doesnt respond
	it should be evicted from the bucket and the contact should be added. */
func TestRoutingTable_1004(t *testing.T) {
	node1 := NewNode ("1111111100000000000000000000000000000000", "localhost:8001");
	node2 := NewNode ("1111111200000000000000000000000000000000", "localhost:8002");
	node3 := NewNode ("1111111300000000000000000000000000000000", "localhost:8003");
	node4 := NewNode ("1111111400000000000000000000000000000000", "localhost:8004");
	node5 := NewNode ("1111111500000000000000000000000000000000", "localhost:8005");
	node6 := NewNode ("1111111600000000000000000000000000000000", "localhost:8006");

	node1.Rt.AddContact(*node2.Me)
	node1.Rt.AddContact(*node3.Me)
	node1.Rt.AddContact(*node4.Me)
	node1.Rt.AddContact(*node5.Me)
	node1.Rt.AddContact(*node6.Me)

	contacts := node1.Rt.FindClosestContacts(node2.Me.ID, 10)
	for i := range contacts {
		log.Println(contacts[i].ID.String())
	}

	/*log.Println("Bucket for " + contact1.ID.String() + ": " + strconv.Itoa(rt.getBucketIndex(contact1.ID)))
	log.Println("Bucket for " + contact2.ID.String() + ": " + strconv.Itoa(rt.getBucketIndex(contact2.ID)))
	log.Println("Bucket for " + contact3.ID.String() + ": " + strconv.Itoa(rt.getBucketIndex(contact3.ID)))
	log.Println("Bucket for " + contact4.ID.String() + ": " + strconv.Itoa(rt.getBucketIndex(contact4.ID))) */
	//log.Println("Bucket for " + node2.Address + strconv.Itoa(rt.getBucketIndex(node1.ID)) + ":")
	/*for i := range rt.buckets {
		log.Println("Bucket " + strconv.Itoa(i) + ":")
		log.Println(i)
	} */


}



