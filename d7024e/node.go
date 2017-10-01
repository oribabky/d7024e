package d7024e

import (
	"log"
)

type Node struct {
	Me *Contact
	Rt *RoutingTable
	Kademlia *Kademlia
	network *Network
}

func NewNode (nodeID string, address string) *Node{

	//see if a nodeID was provided or not.
	me := NewContact(NewRandomKademliaID(), address);
	if nodeID != "" {
		me = NewContact(NewKademliaID(nodeID), address);
	}

	network := NewNetwork(&me)
	rt := NewRoutingTable(me, &network)
	kademlia := NewKademlia(rt, &network)

	return &Node{&me, rt, &kademlia, &network}
}

func (node *Node) NodeUp () {
	log.Println("Hello I am node " + node.Me.ID.String() + " I am on address: " + node.Me.Address)
	go node.network.Listen()
	node.network.RequestHandler(node.Rt)
}