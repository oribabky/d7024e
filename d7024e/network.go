package d7024e
import (
	"log"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"net"
	"sync"
	"time"
	"errors"
	)

type Network struct {
	Contact *Contact
	packetQueue chan *KademliaPacket
	packetID int32
	sentPackets []*KademliaPacket
	mux sync.Mutex
	connection *net.UDPConn
	ReturnedContacts chan *Contact
	files []*File
	ReturnedPacketFiles chan *filePacket   //JOBBA PÅ HÄR!!!!
}

type File struct {
	Key *KademliaID
	Data []byte
}

func NewFile(id string, data []byte) File{

	if id == "" {
		return File{NewRandomKademliaID(), data}
	} 
	return File{NewKademliaID(id), data}
}

func NewNetwork(contact *Contact) Network {
	serverAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	CheckError(err, "resolveError")
	connection, err := net.ListenUDP("udp", serverAddr)
	CheckError(err, "listenError")

	return Network{contact, make(chan *KademliaPacket), 0, make([]*KademliaPacket, 0), sync.Mutex{}, connection, make(chan *Contact), make([]*File, 0), make(chan *File, 1)}
}

//protocol for how rpcs should be written as strings
const PingSend string = "pingSend"
const PingReq string = "pingRequest"
const PingResp string = "pingResponse"

const FindNodeSend string = "findNodeSend"
const FindNodeReq string = "findNodeRequest"
const FindNodeResp string = "findNodeResponse"

const FindDataSend string = "findDataSend"
const FindDataReq string = "findDataRequest"
const FindDataResp string = "findDataResponse"

const StoreSend string = "storeSend"
const StoreReq string = "storeRequest"

func (network *Network) FileExists(fileKey *KademliaID) bool {
/* This function checks if a file already is stored here */
	for i := range network.files {
		if fileKey.String() == network.files[i].Key.String() {
			return true;
		}
	}
	return false;
}

func (network *Network) ReservePacketID(packet *KademliaPacket) int32 {
	/* This function will append a packet to sentPackets[] and incremenet packetID. 
	We need to lock the access to packetID. */
	network.mux.Lock()
	oldValue := network.packetID;
	network.packetID++
	network.sentPackets = append(network.sentPackets, packet)
	packet.PacketID = oldValue;
	defer network.mux.Unlock()
	return oldValue;
}

func (network *Network) RequestHandler(rt *RoutingTable) {
	//Handles requests coming from the packetQueue.
	for {

		currentPacket := <-network.packetQueue
		log.Println("Node " + network.Contact.Address + " handling: " + currentPacket.Procedure + 
			" from " + currentPacket.SourceAddress)

		switch currentPacket.Procedure {

		//PING
		case PingReq:
			kademliaPacket, err := network.CreateKademliaPacket(network.Contact.Address, network.Contact.ID.String(), PingResp)
			CheckError(err, "Error with pingreq")

			kademliaPacket.PacketID = currentPacket.PacketID;
			network.SendKademliaPacket(currentPacket.SourceAddress, kademliaPacket)

		case PingResp:
			log.Println("Pinged and received response from " + 
				currentPacket.SourceAddress)
			network.MarkReturnedPacket(currentPacket)

		case PingSend:
			currentPacket.Procedure = PingReq
			network.SendKademliaPacket(currentPacket.DestinationAddress, currentPacket)


		//FIND_NODE
		case FindNodeReq:
			//add to routing table
			rt.AddContact(NewContact(NewKademliaID(currentPacket.SourceID), currentPacket.SourceAddress))

			targetID := NewKademliaID(currentPacket.TargetID)
			kClosest := rt.FindClosestContacts(targetID, K)

			kademliaPacket, err := network.CreateKademliaPacket(network.Contact.Address, network.Contact.ID.String(), FindNodeResp)
			CheckError(err, "Error with find node req")

			kademliaPacket.PacketID = currentPacket.PacketID;

			for i := range kClosest {
				contactPacket := ContactPacket {
					Address: kClosest[i].Address,
					ID: kClosest[i].ID.String(),
				}
				kademliaPacket.Contacts = append(kademliaPacket.Contacts, &contactPacket)
			} 

			network.SendKademliaPacket(currentPacket.SourceAddress, kademliaPacket)

		case FindNodeResp:
			log.Println("Find_node response received from " + 
				currentPacket.SourceAddress)

						//add to routing table
			rt.AddContact(NewContact(NewKademliaID(currentPacket.SourceID), currentPacket.SourceAddress))

			//network.MarkReturnedPacket(currentPacket)
			for i := range currentPacket.Contacts {
				c := NewContact(NewKademliaID(currentPacket.Contacts[i].ID), currentPacket.Contacts[i].Address)
				go network.AddToContactChannel(&c);
			}
			//for i := range rpc.
			//find k closest nodes to the target ID from my routing table.
		
		case FindNodeSend:
			currentPacket.Procedure = FindNodeReq
			network.SendKademliaPacket(currentPacket.DestinationAddress, currentPacket)


		//FIND_DATA
		case FindDataReq:
			//add to routing table
			rt.AddContact(NewContact(NewKademliaID(currentPacket.SourceID), currentPacket.SourceAddress))

			targetID := NewKademliaID(currentPacket.TargetID)
			kClosest := rt.FindClosestContacts(targetID, K)

			kademliaPacket, err := network.CreateKademliaPacket(network.Contact.Address, network.Contact.ID.String(), FindDataResp)
			CheckError(err, "Error with find data req")
			kademliaPacket.PacketID = currentPacket.PacketID;

			//if the file doesnt exist here, return k-closest contacts to the target ID
			if network.FileExists(targetID) == false {
				for i := range kClosest {
					contactPacket := ContactPacket {
						Address: kClosest[i].Address,
						ID: kClosest[i].ID.String(),
					}
					kademliaPacket.Contacts = append(kademliaPacket.Contacts, &contactPacket)
				} 
			} else {
				//return the file
				data := make([]byte, 0)

				for i := range network.files {
					if network.files[i].Key.String() == targetID.String() {
						data = append(data, network.files[i].Data...)
					}
				}
				filePacket := FilePacket {
					ID: targetID.String(),
					Data: data,
				}
				kademliaPacket.File = &filePacket;
			}


			network.SendKademliaPacket(currentPacket.SourceAddress, kademliaPacket)

		case FindDataResp:
			log.Println("Find_data response received from " + 
				currentPacket.SourceAddress)

			//add to routing table
			rt.AddContact(NewContact(NewKademliaID(currentPacket.SourceID), currentPacket.SourceAddress))

			//if a file has been returned:
			if currentPacket.File.ID != "" {
				log.Println("File retrieved!: " + currentPacket.File.ID)
				file := NewFile(currentPacket.File.ID, currentPacket.File.Data)
				go network.AddToFileChannel(&file)
			} else {
			//if no file was returned, we will return the closest contacts to the key that was returned.
				for i := range currentPacket.Contacts {
					c := NewContact(NewKademliaID(currentPacket.Contacts[i].ID), currentPacket.Contacts[i].Address)
					go network.AddToContactChannel(&c);
				}
			}
			//network.MarkReturnedPacket(currentPacket)

		case FindDataSend:
			currentPacket.Procedure = FindDataReq
			network.SendKademliaPacket(currentPacket.DestinationAddress, currentPacket)



		//STORE
		case StoreReq:
			rt.AddContact(NewContact(NewKademliaID(currentPacket.SourceID), currentPacket.SourceAddress))

			//add the file to the list of files if the file does not already exist here.
			file := NewFile(currentPacket.File.ID, currentPacket.File.Data)
			
			if network.FileExists(file.Key) == false {
				network.files = append(network.files, &file)
				log.Println("Stored file: " + file.Key.String() + " data: " + string(file.Data))
			}

		case StoreSend:
			currentPacket.Procedure = StoreReq
			network.SendKademliaPacket(currentPacket.DestinationAddress, currentPacket)
		}
	}
	
}

func (network *Network) Listen() {
	buf := make([]byte, 1024)

	for {
		//log.Println("listening...")
		n, addr, err := network.connection.ReadFromUDP(buf)
		kademliaPacket := &KademliaPacket{}
		err = proto.Unmarshal(buf[0:n], kademliaPacket)
		if addr != nil {
			go network.AddToPacketChannel(kademliaPacket)
			log.Printf("Received RPC-request: " + kademliaPacket.Procedure + " from " + kademliaPacket.SourceAddress)
		}

		CheckError(err, "Couldn't listen ")
		defer network.connection.Close()
	}
	
}

func (network *Network) CloseConnection() {
	network.connection.Close()
}

func (network *Network) AddToPacketChannel(packet *KademliaPacket) {
	network.packetQueue <- packet;
}

func (network *Network) AddToContactChannel(contact *Contact) {
	network.ReturnedContacts <- contact;
}

func (network *Network) AddToFileChannel(file *File) {
	network.ReturnedFiles <- file;
}

func (network *Network) SendKademliaPacket(address string, packet *KademliaPacket) {
	/* establish a connection to the target server. */

	targetAddr, err := net.ResolveUDPAddr("udp", address)
	CheckError(err, "targetAddr")
	/*localAddr, err := net.ResolveUDPAddr("udp", network.Contact.Address)
	CheckError(err, "localAddr")
	conn, err := net.DialUDP("udp", localAddr, targetAddr)
	CheckError(err, "dialUDP") */

	data, err := proto.Marshal(packet)
	CheckError(err, "Couldn't marshal the message")

	buf := []byte(data)

	_, err = network.connection.WriteToUDP(buf, targetAddr)
	CheckError(err, "Couldn't write the message")

}

func (network *Network) CreateKademliaPacket(sourceAddress string, sourceID string, procedure string) (packet *KademliaPacket, err error) {

	//check that the procedure is one defined by the constants in this file.
	if procedure != PingReq && procedure != PingResp && procedure != FindNodeReq && procedure != FindNodeResp && procedure != PingSend && procedure != FindNodeSend && procedure != StoreSend && procedure != StoreReq && procedure != FindDataSend && procedure != FindDataReq && procedure != FindDataResp{
		return nil, errors.New(" Bad procedure...")
	}

	kademliaPacket := KademliaPacket{
		SourceAddress: sourceAddress,
		SourceID: sourceID,
		Procedure: procedure,
		RandomID: int32(rand.Intn(256)),
	}
	return &kademliaPacket, nil
}

func (network *Network) MarkReturnedPacket (currentPacket *KademliaPacket) {
	currentPacket.ReturnedPacket = true;
	network.sentPackets[currentPacket.PacketID] = currentPacket;
}

func (network *Network) AwaitResponse(packetID int32) bool{
	/* This function will wait for a response from sending a RPC to a node. */

	alive := false;

	start := time.Now()
	limit := 500 * time.Millisecond	//how long time do we wait for a response?
	t := time.Now()
	elapsed := t.Sub(start)

	for network.sentPackets[packetID].ReturnedPacket == false {
		t = time.Now()
		elapsed = t.Sub(start)

		if elapsed > limit {
			break;
		}
	}

	if network.sentPackets[packetID].ReturnedPacket == true {
		alive = true;
		log.Println("Response received!")
	} else {
		log.Println("Time out on waiting for response..")
	}
	
	return alive;
}

func (network *Network) SendPingMessage(address string) bool {
	kademliaPacket, err := network.CreateKademliaPacket(network.Contact.Address, network.Contact.ID.String(), PingSend)
	CheckError(err, "ping failed")

	kademliaPacket.PacketID = network.ReservePacketID(kademliaPacket)
	kademliaPacket.DestinationAddress = address;
	go network.AddToPacketChannel(kademliaPacket)
	//go network.AwaitResponse(kademliaPacket.PacketID)
	alive := network.AwaitResponse(kademliaPacket.PacketID)
	return alive
}

func (network *Network) SendFindNodeMessage(address string, targetID string) {
	kademliaPacket, err := network.CreateKademliaPacket(network.Contact.Address, network.Contact.ID.String(), FindNodeSend)
	CheckError(err, "find_node failed")

	kademliaPacket.PacketID = network.ReservePacketID(kademliaPacket)
	kademliaPacket.TargetID = targetID;

	kademliaPacket.DestinationAddress = address;
	go network.AddToPacketChannel(kademliaPacket)

	//network.AwaitResponse(kademliaPacket.PacketID)
}

func (network *Network) SendFindDataMessage(address string, keyID string) {
	kademliaPacket, err := network.CreateKademliaPacket(network.Contact.Address, network.Contact.ID.String(), FindDataSend)
	CheckError(err, "find_data failed")

	kademliaPacket.PacketID = network.ReservePacketID(kademliaPacket)
	kademliaPacket.TargetID = keyID;

	kademliaPacket.DestinationAddress = address;
	go network.AddToPacketChannel(kademliaPacket)
}

func (network *Network) SendStoreMessage(address string, file *File) {
	kademliaPacket, err := network.CreateKademliaPacket(network.Contact.Address, network.Contact.ID.String(), StoreSend)
	CheckError(err, "store failed")
	kademliaPacket.PacketID = network.ReservePacketID(kademliaPacket)
	kademliaPacket.DestinationAddress = address;

	//file := NewFile("", data)
	//log.Println(file.key.String())

	filePacket := FilePacket {
		ID: file.Key.String(),
		Data: file.Data,
	}

	kademliaPacket.File = &filePacket;

	go network.AddToPacketChannel(kademliaPacket)

}

func CheckError(err error, message string) {
	if err != nil {
		log.Fatal("Error: " + message, err)
	}
}