package d7024e
import (
	"log"
	"github.com/golang/protobuf/proto"
	//"fmt"
	"net"
	)

type RPC struct {
	srcAddress string
	procedure string
	targetID string
}

type Network struct {
	contact *Contact
	channel chan *RPC
	kademlia *Kademlia
}

//protocol for how rpcs should be written as strings
const PingReq string = "pingRequest"
const PingResp string = "pingResponse"
const FindNodeReq string = "findNodeRequest"
const FindNodeResp string = "findNodeResponse"

func NewRPC(srcAddress string, procedure string, targetID string) RPC {
	return RPC{srcAddress, procedure, targetID}
}

func NewNetwork(contact *Contact, kademlia *Kademlia) Network {
	return Network{contact, make(chan *RPC), kademlia}
}

func (network *Network) RequestHandler() {
	//Handles requests coming from the channel.
	for {

		requestProcedure := <-network.channel
		log.Println("handling: " + requestProcedure.procedure + 
			" from " + requestProcedure.srcAddress)

		switch requestProcedure.procedure {
		case PingReq:
			//THIS MIGHT NEED TO BE CHANGED SO SENDKADEMLIAPACKET
			//DOESNT NEED A CONTACT EACH TIME.
			srcNode := "FFFFFFFF00000000000000000000000000000000";
			target := NewContact(NewKademliaID(srcNode), requestProcedure.srcAddress);

			kademliaPacket := network.CreateKademliaPacket(network.contact.Address, PingResp)
			network.SendKademliaPacket(&target, kademliaPacket)

		case PingResp:
			log.Println("Pinged and received response from " + 
				requestProcedure.srcAddress)

		case FindNodeReq:
			log.Println("hey")
			targetID := NewKademliaID(requestProcedure.targetID)
			alphaClosest := network.kademlia.rt.FindClosestContacts(targetID, Alpha)
		
			for i := range alphaClosest {
				log.Println(alphaClosest[i].ID.String())
			} 

			kademliaPacket := network.CreateKademliaPacket(network.contact.Address, FindNodeResp)
			kademliaPacket.Contacts = alphaClosest
			network.SendKademliaPacket(&target, &kademliaPacket)


		case FindNodeResp:
			//find k closest nodes to the target ID from my routing table.
			
		}
	}
	/*switch rpc {
	case PingReq:
		channel <- rpc
	case PingResp:
		channel <- rpc
	default:
		log.Println("unknown RPC")
	} */
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
			rpcRequest := NewRPC(*kademliaPacket.SourceAddress, *kademliaPacket.Procedure, *kademliaPacket.TargetID)
			go network.AddToChannel(&rpcRequest)
			log.Printf("Received RPC-request: " + *kademliaPacket.Procedure + " from " + *kademliaPacket.SourceAddress)
		}

		CheckError(err, "Couldn't listen ")
	}
	
}

func (network *Network) AddToChannel(rpc *RPC) {
	network.channel <- rpc;
}

func (network *Network) SendKademliaPacket(targetNode *Contact, packet *KademliaPacket) {
	
	//establish a connection to the target server.

	targetAddr, err := net.ResolveUDPAddr("udp", targetNode.Address)
	CheckError(err, "")
	localAddr, err := net.ResolveUDPAddr("udp", network.contact.Address)
	CheckError(err, "")
	conn, err := net.DialUDP("udp", localAddr, targetAddr)
	CheckError(err, "")
	defer conn.Close() //if there is an error, close the connection

	

	if targetID != "" {
		kademliaPacket.TargetID = &targetID;
	}

	data, err := proto.Marshal(kademliaPacket)
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
		SourceAddress: &sourceAddress,
		Procedure: &procedure,
	}
	return &kademliaPacket
}

func (network *Network) SendFindContactMessage(contact *Contact, targetID string) {
	network.SendKademliaPacket(contact, FindNodeReq, targetID)
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