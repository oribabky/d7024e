package d7024e

import (
	"log"
)

type Node struct {
	Me *Contact
	Rt *RoutingTable
	Kademlia *Kademlia
	Network *Network
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
	go node.Network.Listen()
	node.Network.RequestHandler(node.Rt)
}