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
}

//protocol for how rpcs should be written as strings
const PingReq string = "pingRequest"
const PingResp string = "pingResponse"
const FindNodeReq string = "findNodeRequest"
const FindNodeResp string = "findNodeResponse"

func NewRPC(srcAddress string, procedure string, targetID string) RPC {
	return RPC{srcAddress, procedure, targetID}
}

func NewNetwork(contact *Contact) Network {
	return Network{contact, make(chan *RPC)}
}

func (network *Network) RequestHandler(rt *RoutingTable) {
	//Handles requests coming from the channel.
	for {

		rpc := <-network.channel
		log.Println("handling: " + rpc.procedure + 
			" from " + rpc.srcAddress)

		switch rpc.procedure {
		case PingReq:
			/*//THIS MIGHT NEED TO BE CHANGED SO SENDKADEMLIAPACKET
			//DOESNT NEED A CONTACT EACH TIME.
			srcNode := "FFFFFFFF00000000000000000000000000000000";
			target := NewContact(NewKademliaID(srcNode), rpc.srcAddress); */

			kademliaPacket := network.CreateKademliaPacket(network.contact.Address, PingResp)
			network.SendKademliaPacket(rpc.srcAddress, kademliaPacket)

		case PingResp:
			log.Println("Pinged and received response from " + 
				rpc.srcAddress)

		case FindNodeReq:
			targetID := NewKademliaID(rpc.targetID)
			kClosest := rt.FindClosestContacts(targetID, K)

			kademliaPacket := network.CreateKademliaPacket(network.contact.Address, FindNodeResp)

			for i := range kClosest {
				log.Println(kClosest[i].ID.String())
				contactPacket := ContactPacket {
					Address: kClosest[i].Address,
					ID: kClosest[i].ID.String(),
				}
				kademliaPacket.Contacts = append(kademliaPacket.Contacts, &contactPacket)
			} 

			network.SendKademliaPacket(rpc.srcAddress, kademliaPacket)


		case FindNodeResp:
			log.Println("Find_node response received from " + 
				rpc.srcAddress)
			//for i := range rpc.
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
			rpcRequest := NewRPC(kademliaPacket.SourceAddress, kademliaPacket.Procedure, kademliaPacket.TargetID)
			go network.AddToChannel(&rpcRequest)
			log.Printf("Received RPC-request: " + kademliaPacket.Procedure + " from " + kademliaPacket.SourceAddress)
		}

		CheckError(err, "Couldn't listen ")
	}
	
}

func (network *Network) AddToChannel(rpc *RPC) {
	network.channel <- rpc;
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

func (network *Network) SendPingMessage(address string) {
	kademliaPacket := network.CreateKademliaPacket(network.contact.Address, PingReq)
	network.SendKademliaPacket(address, kademliaPacket)
}

func (network *Network) SendFindNodeMessage(address string, targetID string) {
	kademliaPacket := network.CreateKademliaPacket(network.contact.Address, FindNodeReq)
	kademliaPacket.TargetID = targetID;
	network.SendKademliaPacket(address, kademliaPacket)
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