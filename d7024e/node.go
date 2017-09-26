package d7024e

import (
	"log"
)

type Node struct {
	Me *Contact
	Rt *RoutingTable
	kademlia *Kademlia
	network *Network
}

func NewNode (nodeID string, address string) *Node{

	me := NewContact(NewRandomKademliaID(), address);
	if nodeID != "" {
		me = NewContact(NewKademliaID(nodeID), address);
	}/* else {
		me := NewContact(NewRandomKademliaID(), address);
	}*/
	
	rt := NewRoutingTable(me)
	kademlia := NewKademlia(&me, rt)
	network := NewNetwork(&me)
	return &Node{&me, rt, &kademlia, &network}
}

func (node *Node) NodeUp () {
	log.Println("Hello I am node " + node.Me.ID.String() + " I am on address: " + node.Me.Address)
	go node.network.Listen()
	node.network.RequestHandler(node.Rt)
}