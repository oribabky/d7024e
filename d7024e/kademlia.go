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



func (kademlia *Kademlia) LookupContact(target *Contact, k int, alpha int) []Contact {
	//contacts := kademlia.rt.FindClosestContacts(target.ID, intk)

	//selected the alpha closest from our own routing table to the target
	myKClosest := kademlia.rt.FindClosestContacts(target.ID, k)

	//WHITEBOXTEST
	/*log.Println("\nCurrentKClosest From RT for node " + kademlia.network.Contact.Address +":")
	for i := range myKClosest {
		log.Println(myKClosest[i].Address)
	} */
	//

	kClosest := myKClosest;
	/*if kClosest[0].ID.String() == target.ID.String() {	//if the targetID is in kClosest:
		kClosest = kClosest[1:]
	}*/

	toBeQueried := kClosest
	if len(kClosest) > alpha {	//if there are more than alpha entries.
		toBeQueried = kClosest[0:alpha]
	}
	
	queriedContacts := make([]Contact, 0)
	//queriedContacts = append(queriedContacts, *kademlia.network.Contact)


	kClosest = kademlia.NodeLookup(&toBeQueried, &kClosest, &queriedContacts, target, k, alpha)
	return kClosest;
}

func (kademlia *Kademlia) NodeLookup(toBeQueried *[]Contact, kClosest *[]Contact, queriedContacts *[]Contact, target *Contact, k int, alpha int) []Contact {
	toBeQueried1 := *toBeQueried;
	kClosest1 := *kClosest;
	queriedContacts1 := *queriedContacts;

	//WHITEBOXTEST
	log.Println("CurrentKClosest for node " + kademlia.network.Contact.Address +":")
	for i := range kClosest1 {
		log.Println(kClosest1[i].Address)
	} 

		log.Println("\nnodes to be queried:")
	PrintContactList(toBeQueried1)
log.Println("\n")
	/*log.Println("\nTOBEQUERIED:")
	for i := range toBeQueried1 {
		log.Println(toBeQueried1[i].Address)
	}*/
	
	//base case
	if len(toBeQueried1) == 0 {
		return kClosest1;
	}

	for i := range toBeQueried1 {
		go kademlia.network.SendFindNodeMessage(toBeQueried1[i].Address, target.ID.String())
		queriedContacts1 = append(queriedContacts1, toBeQueried1[i])

	}

	toBeQueried1 = ClearContactSlice(toBeQueried1)


	//roundSuccessful := false;

	for {
	    select {
	        case <-time.After(time.Millisecond * 600):
		    	log.Println("timeout!!")
		    	break;

	    	case c := <-kademlia.network.ReturnedContacts:

	    	    //check that c is not already in kClosest1.
				if ContainsContact(kClosest1, c) == true {
					log.Println("contact" + c.Address + "already in kClosest1!")
					continue;
				//check that c is not the target itself
				/*} else if c.ID.String() == target.ID.String() {
					log.Println("cannot add the target itself!") */

				//if kClosest1 holds k items in the array.
				} else if len(kClosest1) >= k {  
		    	    //add the contact to k-closest
		    	    kClosest1 = InsertContactSortedDistTarget(c, kClosest1, target)
		    	    
		    	    //if at least one contact was not inserted on the last index, means that it was of closer distance than some other
		    	    //contact in kClosest1 to our target.
		    	    if kClosest1[k].ID.String() != c.ID.String() {
		    	    	//roundSuccessful = true;
		    	    	log.Println("contact " + c.Address + " was added!")
		    	    }


		    	    //and strip the list to K items
		    	    kClosest1 = kClosest1[0:k]
		    	    PrintContactList(kClosest1)
		    	//if kClosest1 holds less than K items
				} else if len(kClosest1) < k {
					
		    	    //add the contact to k-closest
		    	    kClosest1 = InsertContactSortedDistTarget(c, kClosest1, target)

		    	    log.Println("contact " + c.Address + " was added!")
					lastIndex := len(kClosest1) - 1;	
					PrintContactList(kClosest1)
		    	    //if at least one contact was not inserted on the last index, means that it was of closer distance than some other
		    	    //contact in kClosest1 to our target.
		    	    if kClosest1[lastIndex].ID.String() != c.ID.String() {
		    	    	//roundSuccessful = true;
		    	    }

				}
				continue;		//go back to the switch case.
				
			
			}
		break;	//break out of the outer for-loop.
	}
	/*log.Println("hej")
	PrintContactList(kClosest1)
	log.Println(roundSuccessful) */

	limit := alpha
	/*if roundSuccessful == true {		//pick alpha contacts from kClosest1 that have not yet been queried if the round was successful,
		limit = alpha                  //otherwise we will pick all from kClosest1 that have not been queried.
	}*/

	contactsToQuery := 0
	for i := range kClosest1 {
		alreadyQueried := false;
		currentContact := &kClosest1[i];

		if contactsToQuery >= limit {
			break;
		}

		if ContainsContact(queriedContacts1, currentContact) == true {
			alreadyQueried = true;
		}

		if alreadyQueried == false{
			toBeQueried1 = append(toBeQueried1, *currentContact)
			contactsToQuery ++;
		}
		//log.Println(kClosest1[i].Address)
	}

/*	log.Println("\nafter break")
	PrintContactList(kClosest1)	
	log.Println("\nnodes to be queried:")
	PrintContactList(toBeQueried1) */
	

	/* for i := range kClosest1 {
		if kClosest1[i].ID.String() != oldList[i].ID.String() {
			log.Println("\n\n\nWTFFFFFFFFFFFFFFFFFFf\n\n\n")
		}
	} */
	return kademlia.NodeLookup(&toBeQueried1, &kClosest1, &queriedContacts1, target, k, alpha)
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
		} /*else if contact.Less(currentContact) == currentContact.Less(contact) {		//if the distance are the same to one another.
			index = i;
			break; 
		}*/
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