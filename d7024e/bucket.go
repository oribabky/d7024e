package d7024e

import (
	"container/list"
	"log"
)

type bucket struct {
	list *list.List
}

func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	return bucket
}

func (bucket *bucket) AddContact(contact Contact, network *Network) {
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}

	// if the element doesnt already exists in our list
	if element == nil {
		if bucket.list.Len() < bucketSize {	//add to the bucket
			bucket.list.PushFront(contact)	
		} else {	//ping the least recently seen item and see if its still alive
			log.Println("bucket full! " + "I am " + contact.ID.String() + " pinging LRS contact..")
			leastRecentlySeen := bucket.list.Back().Value.(Contact)
			log.Println("LSR: " + leastRecentlySeen.Address)
			alive := network.SendPingMessage(leastRecentlySeen.Address)

			if alive == false {		//remove the least recently seen item and add the name item
				bucket.list.Remove(bucket.list.Back())
				bucket.list.PushFront(contact)
			} else {	//if the LSR node responds, move it to the front of the list.
				bucket.list.MoveToFront(bucket.list.Back())
			}
		}
	} else {
		bucket.list.MoveToFront(element)
	}
}

func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts
}

func (bucket *bucket) Len() int {
	return bucket.list.Len()
}
