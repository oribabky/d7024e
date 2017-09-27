package d7024e
import (
	"log"
	"github.com/golang/protobuf/proto"
	//"fmt"
	"net"
	"sync"
	)

/*type RPC struct {
	srcAddress string
	procedure string
	targetID string
}*/

type Network struct {
	contact *Contact
	channel chan *KademliaPacket
	packetID int32
	sentPackets []*KademliaPacket
	mux sync.Mutex
}


//protocol for how rpcs should be written as strings
const PingReq string = "pingRequest"
const PingResp string = "pingResponse"
const FindNodeReq string = "findNodeRequest"
const FindNodeResp string = "findNodeResponse"

/*func NewRPC(srcAddress string, procedure string, targetID string) RPC {
	return RPC{srcAddress, procedure, targetID}
}*/

func NewNetwork(contact *Contact) Network {
	return Network{contact, make(chan *KademliaPacket), 0, make([]*KademliaPacket, 0), sync.Mutex{}}
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
	//Handles requests coming from the channel.
	for {

		currentPacket := <-network.channel
		log.Println("handling: " + currentPacket.Procedure + 
			" from " + currentPacket.SourceAddress)

		switch currentPacket.Procedure {
		case PingReq:
			kademliaPacket := network.CreateKademliaPacket(network.contact.Address, PingResp)
			kademliaPacket.PacketID = currentPacket.PacketID;

			network.SendKademliaPacket(currentPacket.SourceAddress, kademliaPacket)

		case PingResp:
			log.Println("Pinged and received response from " + 
				currentPacket.SourceAddress)
			network.MarkReturnedPacket(currentPacket)

		case FindNodeReq:
			targetID := NewKademliaID(currentPacket.TargetID)
			kClosest := rt.FindClosestContacts(targetID, K)

			kademliaPacket := network.CreateKademliaPacket(network.contact.Address, FindNodeResp)
			kademliaPacket.PacketID = currentPacket.PacketID;

			for i := range kClosest {
				log.Println(kClosest[i].ID.String())
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
			network.MarkReturnedPacket(currentPacket)
			for i := range currentPacket.Contacts {
				log.Println(currentPacket.Contacts[i].ID)
			}
			//for i := range rpc.
			//find k closest nodes to the target ID from my routing table.
			
		}
	}
	
}

func (network *Network) Listen() {
	buf := make([]byte, 1024)

	//establish a connection 
	serverAddr, err := net.ResolveUDPAddr("udp", network.contact.Address)
	CheckError(err, "")
	serverConn, err := net.ListenUDP("udp", serverAddr)
	CheckError(err, "")
	defer serverConn.Close() //close the connection when something is return

	for {
		log.Println("listening...")
		n, addr, err := serverConn.ReadFromUDP(buf)
		kademliaPacket := &KademliaPacket{}
		err = proto.Unmarshal(buf[0:n], kademliaPacket)
		if addr != nil {
			//rpcRequest := NewRPC(kademliaPacket.SourceAddress, kademliaPacket.Procedure, kademliaPacket.TargetID)
			go network.AddToChannel(kademliaPacket)
			log.Printf("Received RPC-request: " + kademliaPacket.Procedure + " from " + kademliaPacket.SourceAddress)
		}

		CheckError(err, "Couldn't listen ")
	}
	
}

func (network *Network) AddToChannel(packet *KademliaPacket) {
	network.channel <- packet;
}

func (network *Network) SendKademliaPacket(address string, packet *KademliaPacket) {
	
	//establish a connection to the target server.

	targetAddr, err := net.ResolveUDPAddr("udp", address)
	CheckError(err, "")
	localAddr, err := net.ResolveUDPAddr("udp", network.contact.Address)
	CheckError(err, "")
	conn, err := net.DialUDP("udp", localAddr, targetAddr)
	CheckError(err, "")
	defer conn.Close() //if there is an error, close the connection

	data, err := proto.Marshal(packet)
	CheckError(err, "Couldn't marshal the message")

	buf := []byte(data)

	_, err = conn.Write(buf)
	CheckError(err, "Couldn't write the message")

}

func (network *Network) CreateKademliaPacket(sourceAddress string, procedure string) *KademliaPacket {

	//check that the procedure is one defined by the constants in this file.
	if procedure != PingReq && procedure != PingResp && procedure != FindNodeReq && procedure != FindNodeResp {
		log.Println("bad procedure.." + procedure) //NEED ERROR HANDLING
	}

	kademliaPacket := KademliaPacket{
		SourceAddress: sourceAddress,
		Procedure: procedure,
	}


	return &kademliaPacket
}
func (network *Network) MarkReturnedPacket (currentPacket *KademliaPacket) {
	currentPacket.ReturnedPacket = true;
	network.sentPackets[currentPacket.PacketID] = currentPacket;
}


func (network *Network) AwaitResponse(packetID int32) {
	/* This function will wait for a response from sending a RPC to a node. */

	for network.sentPackets[packetID].ReturnedPacket == false {
		//log.Println(network.sentPackets[packetID].PacketID)
	}
	log.Println("Response received!")
}

func (network *Network) SendPingMessage(address string) {
	kademliaPacket := network.CreateKademliaPacket(network.contact.Address, PingReq)

	reservedID := network.ReservePacketID(kademliaPacket)
	kademliaPacket.PacketID = reservedID;

	network.SendKademliaPacket(address, kademliaPacket)
	network.AwaitResponse(kademliaPacket.PacketID)
}

func (network *Network) SendFindNodeMessage(address string, targetID string) {
	kademliaPacket := network.CreateKademliaPacket(network.contact.Address, FindNodeReq)
	kademliaPacket.TargetID = targetID;

	reservedID := network.ReservePacketID(kademliaPacket)
	kademliaPacket.PacketID = reservedID;

	network.SendKademliaPacket(address, kademliaPacket)
	network.AwaitResponse(kademliaPacket.PacketID)
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func CheckError(err error, message string) {
	if err != nil {
		log.Fatal("Error: " + message, err)
	}
}