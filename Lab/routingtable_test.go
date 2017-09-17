package d7024e

import (
	"testing"
)

/* Test case 1001: FindClosestContacts should return k closest nodes ordered in distance to source node. */
func TestRoutingTable_1001(t *testing.T) {
	

	srcNode := "FFFFFFFF00000000000000000000000000000000";
	node1 := "1111111100000000000000000000000000000000";
	node2 := "1111111200000000000000000000000000000000";
	node3 := "1111111300000000000000000000000000000000";
	node4 := "1111111400000000000000000000000000000000";
	node5 := "2111111400000000000000000000000000000000";

	//expected k closest in the order they should they appear.
	expected := [5]string{node5, node4, node3, node2, node1}

	rt := NewRoutingTable(NewContact(NewKademliaID(srcNode), "localhost:8000"))
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

/* Test case 1002: FindClosestContact should return k items if k is higher than the amount of items available. */
func TestRoutingTable_1002(t *testing.T) {

	node1 := "1111111100000000000000000000000000000000";
	node2 := "1111111200000000000000000000000000000000";
	node3 := "1111111300000000000000000000000000000000";

	rt := NewRoutingTable(NewContact(NewKademliaID(node1), "localhost:8000"))
	rt.AddContact(NewContact(NewKademliaID(node2), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID(node3), "localhost:8002"))

	k := 1000;
	kClosest := rt.FindClosestContacts(NewKademliaID(node1), k)

	expectedSize := 2;

	if len(kClosest) != expectedSize {
			t.Error("error in test case 1002")
		}
}

/* Test case 1003: AddContact should add a contact to the routing table if that ID is not already taken, regardless of port used.*/
func TestRoutingTable_1003(t *testing.T) {
	srcNode := "FFFFFFFF00000000000000000000000000000000";
	node1 := "1111111100000000000000000000000000000000";
	node2 := "1111111200000000000000000000000000000000";

	rt := NewRoutingTable(NewContact(NewKademliaID(srcNode), "localhost:8000"))
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


