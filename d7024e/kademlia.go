package d7024e

import (
//	"fmt"
)

type Kademlia struct {
	rt *RoutingTable
}

func NewKademlia (me *Contact, rt *RoutingTable) Kademlia {
	return Kademlia{rt}
}

const Alpha int = 2;
const K int = 3;



func (kademlia *Kademlia) LookupContact(target *Contact, node *Node) {
	//contacts := kademlia.rt.FindClosestContacts(target.ID, intk)

}


func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
