package d7024e

import (
	"testing"
	"log"
	"time"
	//"fmt"
)

/* Test case 2001: The system should be able to send various RPCs.*/
func Test_2001(t *testing.T) {
	time.Sleep(time.Millisecond * 500)
	log.Println("TEST RPCS..")

	node1 := NewNode("", ClientAddress)
	server1 := NewNode("", ServerAddress1)
	server2 := NewNode("", ServerAddress2)
	server3 := NewNode("", ServerAddress3)
	go node1.NodeUp()
	go server1.NodeUp()
	go server2.NodeUp()
	go server3.NodeUp()

	//TEST PING
	log.Println("\nPING")

	if node1.Network.SendPingMessage(ServerAddress1) == false {
		t.Fatal("error in testing RPCs.")
	}
	if node1.Network.SendPingMessage(ServerAddress2) == false {
		t.Fatal("error in testing RPCs.")
	}
	if node1.Network.SendPingMessage(ServerAddress3) == false {
		t.Fatal("error in testing RPCs.")
	}


	//TEST FIND_NODE, should be able to be sent asynchronously.
	time.Sleep(time.Millisecond * 500)
	log.Println("\nFIND_NODE")

	target := NewContact(NewRandomKademliaID(), "localhost:8000")

	
	//we also need to expect that node1 will be in their routing tables because we send requests to them from node1
	server1.Rt.AddContact(*node1.Me)
	server2.Rt.AddContact(*node1.Me)
	server3.Rt.AddContact(*node1.Me)
	node1.Rt.AddContact(*server1.Me)
	node1.Rt.AddContact(*server2.Me)
	node1.Rt.AddContact(*server3.Me) 


	//send out the find_node rpcs asynchronously
	go node1.Network.SendFindNodeMessage(ServerAddress1, target.ID.String())
	go node1.Network.SendFindNodeMessage(ServerAddress2, target.ID.String())
	go node1.Network.SendFindNodeMessage(ServerAddress3, target.ID.String())
	go node1.Network.SendFindNodeMessage(ServerAddress1, target.ID.String())
	go node1.Network.SendFindNodeMessage(ServerAddress2, target.ID.String())
	go node1.Network.SendFindNodeMessage(ServerAddress3, target.ID.String())
	
	//fetch the actual returned contacts
	kClosestTotalActual := make([]Contact, 0)

	for {
		select {
	        case <-time.After(time.Millisecond * 2000):
		    	log.Println("Channel empty.")
		    	break;

	    	case c := <-node1.Network.ReturnedContacts:
	    		kClosestTotalActual = append(kClosestTotalActual, *c)
	    		continue;
	    	}
	    break;
	}

	//fetch kClosest for each node to the target ID
	kClosestExpected1 := server1.Rt.FindClosestContacts(target.ID, K)
	kClosestExpected2 := server2.Rt.FindClosestContacts(target.ID, K)
	kClosestExpected3 := server3.Rt.FindClosestContacts(target.ID, K)

	kClosestTotalExpected := make([]Contact, 0)
	kClosestTotalExpected = append(kClosestTotalExpected, kClosestExpected1...)
	kClosestTotalExpected = append(kClosestTotalExpected, kClosestExpected2...)
	kClosestTotalExpected = append(kClosestTotalExpected, kClosestExpected3...)
	kClosestTotalExpected = append(kClosestTotalExpected, kClosestExpected1...)
	kClosestTotalExpected = append(kClosestTotalExpected, kClosestExpected2...)
	kClosestTotalExpected = append(kClosestTotalExpected, kClosestExpected3...)


	//check that the size is the same of both slices
	if len(kClosestTotalActual) != len(kClosestTotalExpected) {
		log.Println(len(kClosestTotalActual))
		log.Println(len(kClosestTotalExpected))
		t.Fatal("error in testing RPCs. ")
	}

	//check that the contents are the same. That is, check that the returned contacts match the kClosestTotal:
	 for i := range kClosestTotalExpected {
		currentContact := kClosestTotalExpected[i]
		foundMatch := ContainsContact(kClosestTotalActual, &currentContact)
		if foundMatch == false {
			t.Fatal("error in testing RPCs.")
		}
	} 

	time.Sleep(time.Millisecond * 500)


	//TEST STORE
	log.Println("\nSTORE")

	fileContents := []byte("asdasdasdasdasd")
	file := NewFile("", fileContents)
	node1.Network.SendStoreMessage(ServerAddress1, &file)
	node1.Network.SendStoreMessage(ServerAddress2, &file)
	node1.Network.SendStoreMessage(ServerAddress3, &file)

	time.Sleep(time.Millisecond * 500)

	if server1.Network.FileExists(file.Key) == false || server2.Network.FileExists(file.Key) == false || server3.Network.FileExists(file.Key) == false {
		log.Println("file " + file.Key.String() + " does not exist in one of the servers")
		t.Fatal("error in testing RPCs.")
	} 

	time.Sleep(time.Millisecond * 500)


	//TEST FIND_VALUE
	log.Println("\nFIND_VALUE")
	go node1.Network.SendFindDataMessage(ServerAddress1, file.Key.String())
	go node1.Network.SendFindDataMessage(ServerAddress2, file.Key.String())
	go node1.Network.SendFindDataMessage(ServerAddress3, file.Key.String())


	returnedFileID := "";
	returnedFileData := "";
	for {
		select {
	        case <-time.After(time.Millisecond * 1000):
		    	log.Println("Channel empty.")
		    	break;

	    	case filePacket := <-node1.Network.ReturnedPacketFiles:
	    		log.Println("extracting file: " + filePacket.ID)
	    		log.Println("contents: " + string(file.Data)	)
	    		returnedFileID = filePacket.ID;
	    		returnedFileData = string(file.Data)
	    		break;
	    	}
	    break;
	}

	if returnedFileID == "" || returnedFileData == "" {
		t.Fatal("error in testing RPCs.")
	}


	time.Sleep(time.Millisecond * 500)
	node1.Network.CloseConnection();
	server1.Network.CloseConnection();
	server2.Network.CloseConnection();
	server3.Network.CloseConnection(); 
	time.Sleep(time.Millisecond * 500)
} 