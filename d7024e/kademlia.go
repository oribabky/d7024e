package d7024e

import (
	"fmt"
	"time"
)

type Kademlia struct {
	rt *RoutingTable
	network *Network
}

func NewKademlia (rt *RoutingTable, network *Network) Kademlia {
	return Kademlia{rt, network}
}

const Alpha int = 2;
const K int = 3;



func (kademlia *Kademlia) LookupContact(target *Contact, k int) []Contact {
	//contacts := kademlia.rt.FindClosestContacts(target.ID, intk)

	//selected the alpha closest from our own routing table to the target
	myKClosest := kademlia.rt.FindClosestContacts(target.ID, k)
	kClosest := myKClosest;

	toBeQueried := kClosest
	if len(kClosest) > Alpha {	//if there are more than alpha entries.
		toBeQueried = kClosest[0:Alpha]
	}
	
	queriedContacts := make([]Contact, 0)
	queriedContacts = append(queriedContacts, *kademlia.network.Contact)

	kClosest = kademlia.NodeLookup(toBeQueried, kClosest, queriedContacts, target, k)
	return kClosest;
}

func (kademlia *Kademlia) NodeLookup(toBeQueried []Contact, kClosest []Contact, queriedContacts []Contact, target *Contact, k int) []Contact {
	//WHITEBOXTEST
	/*fmt.Println("\nCurrentKClosest:")
	for i := range kClosest {
		fmt.Println(kClosest[i].ID.String())
	} */
	/*fmt.Println("\nTOBEQUERIED:")
	for i := range toBeQueried {
		fmt.Println(toBeQueried[i].Address)
	}*/
	
	//base case
	if len(toBeQueried) == 0 {
		return kClosest;
	}

	for i := range toBeQueried {
		go kademlia.network.SendFindNodeMessage(toBeQueried[i].Address, target.ID.String())
		queriedContacts = append(queriedContacts, toBeQueried[i])
	}

	toBeQueried = ClearContactSlice(toBeQueried)


	roundSuccessful := false;

	for {
	    select {
	        case <-time.After(time.Millisecond * 1000):
		    	fmt.Println("timeout!!")
		    	break;

	    	case c := <-kademlia.network.ReturnedContacts:
	    	    //check that c is not already in kClosest.
				if ContainsContact(kClosest, *c) == true {
					fmt.Println("contact already in kClosest!")
					continue;
				}

				//check that c is not the target itself
				if c.ID.String() == target.ID.String() {
					fmt.Println("cannot add ourselves!")
					continue;
				}

				if len(kClosest) == k {  //if kClosest holds k items in the array.
		    	    //add the contact to k-closest
		    	    kClosest = InsertContactSortedDistTarget(c, kClosest, target)

		    	    //if at least one contact was not inserted on the last index, means that it was of closer distance than some other
		    	    //contact in kClosest to our target.
		    	    if kClosest[k].ID.String() != c.ID.String() {
		    	    	roundSuccessful = true;
		    	    }
		    	    //and strip the list to K items
		    	    kClosest = kClosest[0:k]

				} else if len(kClosest) < k {
					//if kClosest holds less than K items

		    	    //add the contact to k-closest
		    	    kClosest = InsertContactSortedDistTarget(c, kClosest, target)

					lastIndex := len(kClosest) - 1;	

		    	    //if at least one contact was not inserted on the last index, means that it was of closer distance than some other
		    	    //contact in kClosest to our target.
		    	    if kClosest[lastIndex].ID.String() != c.ID.String() {
		    	    	roundSuccessful = true;
		    	    }
		    		continue;
				}
			}
		break;	//break out of the outer for-loop.
	}

	limit := k
	if roundSuccessful == true {		//pick Alpha contacts from kClosest that have not yet been queried if the round was successful,
		limit = Alpha                  //otherwise we will pick all from kClosest that have not been queried.
	}

	contactsToQuery := 0
	for i := range kClosest {
		if contactsToQuery == limit {
			break;
		}

		alreadyQueried := false;
		for k := range queriedContacts {
			if kClosest[i].ID.String() == queriedContacts[k].ID.String() {
				alreadyQueried = true;
				break;
			}
		}

		if alreadyQueried == false {
			toBeQueried = append(toBeQueried, kClosest[i])
			contactsToQuery ++;
		}
	}

	
	return kademlia.NodeLookup(toBeQueried, kClosest, queriedContacts, target, k)
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

func ContainsContact(contacts []Contact, contact Contact) bool {
    for i := range contacts {
        if contact.ID.String() == contacts[i].ID.String() {
            return true
        }
    }
    return false
}

func InsertContactSortedDistTarget(contact *Contact, list []Contact, target *Contact) []Contact {
	/* This function will insert a contact in a list, the list is sorted on distance to target. */

	//find the right index
	index := len(list)	//initialize it as the last index.
	contact.CalcDistance(target.ID)

	for i := range list {
		currentContact := &list[i]
		currentContact.CalcDistance(target.ID)

		if contact.Less(currentContact) {		//kClosest is sorted on distance to target node.
			index = i;
			break;
		}
	}

	s := append(list, list[0])
	copy(s[index+1:], s[index:])
	s[index] = *contact
	return s
}

func ClearContactSlice(list []Contact) []Contact {
	for len(list) >= 1 {
		list = list[1:]
	}
	return list;
}