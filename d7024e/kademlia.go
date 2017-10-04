package d7024e

import (
	"log"
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



func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	//contacts := kademlia.rt.FindClosestContacts(target.ID, intk)

	//selected the alpha closest from our own routing table to the target
	myKClosest := kademlia.rt.FindClosestContacts(target.ID, K)

	kClosest := make([]Contact, 0)
	kClosest = append(kClosest, myKClosest...)

	toBeQueried := make([]Contact, 0)
	toBeQueried = append(toBeQueried, kClosest...)

	if len(kClosest) > Alpha {	//if there are more than alpha entries.
		toBeQueried = append(toBeQueried, kClosest[0:Alpha]...)
	}
	
	queriedContacts := make([]Contact, 0)

	kClosest = kademlia.NodeLookup(toBeQueried, kClosest, queriedContacts, target)
	return kClosest;
}

func (kademlia *Kademlia) NodeLookup(toBeQueried []Contact, kClosest []Contact, queriedContacts []Contact, target *Contact) []Contact {

	/*//WHITEBOXTEST
	log.Println("CurrentKClosest for node " + kademlia.network.Contact.Address +":")
	PrintContactList(kClosest)

		log.Println("\nnodes to be queried:")
	PrintContactList(toBeQueried)
	log.Println("\n")
	log.Println("\nTOBEQUERIED:")
	for i := range toBeQueried {
		log.Println(toBeQueried[i].Address)
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


	//roundSuccessful := false;
	currentKClosest := kClosest
	for {
	    select {
	        case <-time.After(time.Millisecond * 400):
		    	log.Println("timeout!!")
		    	break;

	    	case c := <-kademlia.network.ReturnedContacts:

	    	    //check that c is not already in currentKClosest.
				if ContainsContact(currentKClosest, c) == true {
					log.Println("contact" + c.Address + "already in currentKClosest!")
					continue;

				//if currentKClosest holds k items in the array.
				} else if len(currentKClosest) >= K {  
		    	    //add the contact to k-closest
		    	    currentKClosest = InsertContactSortedDistTarget(c, currentKClosest, target)
		    	    
		    	    //if at least one contact was not inserted on the last index, means that it was of closer distance than some other
		    	    //contact in currentKClosest to our target.
		    	    if currentKClosest[K].ID.String() != c.ID.String() {
		    	    	//roundSuccessful = true;
		    	    	log.Println("contact " + c.Address + " was added!")
		    	    }


		    	    //and strip the list to K items
		    	    currentKClosest = currentKClosest[0:K]

		    	//if currentKClosest holds less than K items
				} else if len(currentKClosest) < K {
					
		    	    //add the contact to k-closest
		    	    currentKClosest = InsertContactSortedDistTarget(c, currentKClosest, target)

		    	    log.Println("contact " + c.Address + " was added!")
					lastIndex := len(currentKClosest) - 1;	

		    	    //if at least one contact was not inserted on the last index, means that it was of closer distance than some other
		    	    //contact in currentKClosest to our target.
		    	    if currentKClosest[lastIndex].ID.String() != c.ID.String() {
		    	    	//roundSuccessful = true;
		    	    }

				}
				continue;		//go back to the switch case.
				
			
			}
		break;	//break out of the outer for-loop.
	}

	limit := Alpha
	/*if roundSuccessful == true {		//pick Alpha contacts from currentKClosest that have not yet been queried if the round was successful,
		limit = Alpha                  //otherwise we will pick all from currentKClosest that have not been queried.
	}*/

	contactsToQuery := 0
	for i := range currentKClosest {
		alreadyQueried := false;
		currentContact := currentKClosest[i];

		if contactsToQuery >= limit {
			break;
		}

		if ContainsContact(queriedContacts, &currentContact) == true {
			alreadyQueried = true;
		}

		if alreadyQueried == false{
			contactToBeAdded := NewContact(NewKademliaID(currentContact.ID.String()), currentContact.Address)
			toBeQueried = append(toBeQueried, contactToBeAdded)
			contactsToQuery ++;
		}
		
	}

	return kademlia.NodeLookup(toBeQueried, currentKClosest, queriedContacts, target)
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

func ContainsContact(contacts []Contact, contact *Contact) bool {
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
/*	for len(list) >= 1 {
		list = list[1:]
	} 
	return list; */
	return list[:0]
	
}