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
const K int = 5;



func (kademlia *Kademlia) LookupContact(targetID *KademliaID) []Contact {

	//selected the alpha closest from our own routing table to the target
	myKClosest := kademlia.rt.FindClosestContacts(targetID, K)

	kClosest := make([]Contact, 0)
	kClosest = append(kClosest, myKClosest...)

	toBeQueried := make([]Contact, 0)
	toBeQueried = append(toBeQueried, kClosest...)

	if len(kClosest) > Alpha {	//if there are more than alpha entries.
		toBeQueried = append(toBeQueried, kClosest[0:Alpha]...)
	}
	
	queriedContacts := make([]Contact, 0)

	kClosest = kademlia.NodeLookup(toBeQueried, kClosest, queriedContacts, targetID)
	return kClosest;
}

func (kademlia *Kademlia) NodeLookup(toBeQueried []Contact, kClosest []Contact, queriedContacts []Contact, targetID *KademliaID) []Contact {

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
		go kademlia.network.SendFindNodeMessage(toBeQueried[i].Address, targetID.String())
		queriedContacts = append(queriedContacts, toBeQueried[i])

	}

	toBeQueried = ClearContactSlice(toBeQueried)

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

				//if currentKClosest holds k items in the array add the contact to k-closest
				} else if len(currentKClosest) >= K {  
		    	    currentKClosest = InsertContactSortedDistTarget(c, currentKClosest, targetID)
		    	    
		    	    //if at least one contact was not inserted on the last index, means that it was of closer distance than some other contact in currentKClosest to our target.
		    	    if currentKClosest[K].ID.String() != c.ID.String() {
		    	    	//roundSuccessful = true;
		    	    	log.Println("contact " + c.Address + " was added!")
		    	    }

		    	    //and strip the list to K items
		    	    currentKClosest = currentKClosest[0:K]

		    	//if currentKClosest holds less than K items
				} else if len(currentKClosest) < K {
					
		    	    //add the contact to k-closest
		    	    currentKClosest = InsertContactSortedDistTarget(c, currentKClosest, targetID)

		    	    log.Println("contact " + c.Address + " was added!")
					lastIndex := len(currentKClosest) - 1;	

		    	    //if at least one contact was not inserted on the last index, means that it was of closer distance than some other
		    	    //contact in currentKClosest to our target.
		    	    if currentKClosest[lastIndex].ID.String() != c.ID.String() {
		    	    }

				}
				continue;		//go back to the select case.
			}
		break;	//break out of the outer for-loop.
	}

	limit := Alpha

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

	return kademlia.NodeLookup(toBeQueried, currentKClosest, queriedContacts, targetID)
}


func (kademlia *Kademlia) LookupData(targetKey *KademliaID) []byte {

	//selected the alpha closest from our own routing table to the target
	myKClosest := kademlia.rt.FindClosestContacts(targetKey, K)

	kClosest := make([]Contact, 0)
	kClosest = append(kClosest, myKClosest...)

	toBeQueried := make([]Contact, 0)
	toBeQueried = append(toBeQueried, kClosest...)

	if len(kClosest) > Alpha {	//if there are more than alpha entries.
		toBeQueried = append(toBeQueried, kClosest[0:Alpha]...)
	}
	
	queriedContacts := make([]Contact, 0)

	data := kademlia.ValueLookup(toBeQueried, kClosest, queriedContacts, targetKey)
	return data;
} 

func (kademlia *Kademlia) ValueLookup(toBeQueried []Contact, kClosest []Contact, queriedContacts []Contact, targetID *KademliaID) []byte {

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
	
	//base case, when no contacts are left to query we return an empty result.
	if len(toBeQueried) == 0 {
		return nil;
	}

	for i := range toBeQueried {
		go kademlia.network.SendFindDataMessage(toBeQueried[i].Address, targetID.String())		//we will send find data messages instead of find node
		queriedContacts = append(queriedContacts, toBeQueried[i])
	}

	currentKClosest := kClosest
	for {
	    select {
	        case <-time.After(time.Millisecond * 400):
		    	log.Println("timeout!!")
		    	break;

	        case file := <-kademlia.network.ReturnedFiles:
	        	//if the file is returned
		    	log.Println("File returned: " + file.Key.String())

		    	//find the closest observed contact that did not return the file:
		    	for i := range kClosest {

		    		

		    	}

		    	return file.Data;

	    	case c := <-kademlia.network.ReturnedContacts:

	    	    //check that c is not already in currentKClosest.
				if ContainsContact(currentKClosest, c) == true {
					log.Println("contact" + c.Address + "already in currentKClosest!")
					continue;

				//if currentKClosest holds k items in the array add the contact to k-closest
				} else if len(currentKClosest) >= K {  
		    	    currentKClosest = InsertContactSortedDistTarget(c, currentKClosest, targetID)
		    	    
		    	    //if at least one contact was not inserted on the last index, means that it was of closer distance than some other contact in currentKClosest to our target.
		    	    if currentKClosest[K].ID.String() != c.ID.String() {
		    	    	//roundSuccessful = true;
		    	    	log.Println("contact " + c.Address + " was added!")
		    	    }

		    	    //and strip the list to K items
		    	    currentKClosest = currentKClosest[0:K]

		    	//if currentKClosest holds less than K items
				} else if len(currentKClosest) < K {
					
		    	    //add the contact to k-closest
		    	    currentKClosest = InsertContactSortedDistTarget(c, currentKClosest, targetID)

		    	    log.Println("contact " + c.Address + " was added!")
					lastIndex := len(currentKClosest) - 1;	

		    	    //if at least one contact was not inserted on the last index, means that it was of closer distance than some other
		    	    //contact in currentKClosest to our target.
		    	    if currentKClosest[lastIndex].ID.String() != c.ID.String() {
		    	    }

				}
				continue;		//go back to the select case.
			}
		break;	//break out of the outer for-loop.
	}

	toBeQueried = ClearContactSlice(toBeQueried)

	limit := Alpha

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

	return kademlia.ValueLookup(toBeQueried, currentKClosest, queriedContacts, targetID)
}


func (kademlia *Kademlia) Store(data []byte) *KademliaID{
	fileToBeAdded := NewFile("", data)

	//find the closest contacts to this dummy-target contact
	kClosest := kademlia.LookupContact(fileToBeAdded.Key)
	PrintContactList(kClosest)
	//send them store RPCs
	for i := range kClosest {
		kademlia.network.SendStoreMessage(kClosest[i].Address, &fileToBeAdded)
	}

	return fileToBeAdded.Key;
}

func ContainsContact(contacts []Contact, contact *Contact) bool {
    for i := range contacts {
        if contact.ID.String() == contacts[i].ID.String() {
            return true
        }
    }
    return false
}

func InsertContactSortedDistTarget(contact *Contact, list []Contact, target *KademliaID) []Contact {
	/* This function will insert a contact in a list, the list is sorted on distance to target. */

	//find the right index
	index := len(list)	//initialize it as the last index.
	distanceNewContact := contact.ID.CalcDistance(target)

	for i := range list {
		currentContact := &list[i]
		distanceCurrentContact := currentContact.ID.CalcDistance(target)
		//currentContact.CalcDistance(target)

		if distanceNewContact.Less(distanceCurrentContact) {		//kClosest is sorted on distance to target node.
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