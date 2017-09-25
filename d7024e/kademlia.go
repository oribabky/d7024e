package d7024e

import (
//	"fmt"
)

type Kademlia struct {
	//network Network
	rt *RoutingTable
}

func NewKademlia (me *Contact, rt *RoutingTable) Kademlia {
	return Kademlia{rt}
}

const Alpha int = 1;
const intk int = 20;



func (kademlia *Kademlia) LookupContact(target *Contact) {
	//contacts := kademlia.rt.FindClosestContacts(target.ID, intk)


}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
