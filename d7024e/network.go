package d7024e
import (
//	"time"
	"log"
	"github.com/golang/protobuf/proto"
	//"fmt"
	"net"
	)

type RPC struct {
	srcAddress string
	procedure string
}

type Network struct {
	contact *Contact
	channel chan *RPC
}

//protocol for how rpcs should be written as strings
const PingReq string = "pingRequest"
const PingResp string = "pingResponse"


func NewRPC(srcAddress string, procedure string) RPC {
	return RPC{srcAddress, procedure}
}

func NewNetwork(contact *Contact) Network {
	return Network{contact, make(chan *RPC)}
}

func (network *Network) RequestHandler() {
	//Handles requests coming from the channel.
	for {
		log.Println("Requesthandler")
		/*request := <-network.channel
		v := <-ch */

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
			rpcRequest := NewRPC(kademliaPacket.SourceAddress, kademliaPacket.Procedure)
			network.channel <- &rpcRequest;
			log.Printf("Received RPC-request: " + kademliaPacket.Procedure + " from " + kademliaPacket.SourceAddress)
			//log.Printf("Received %s at %s from %s", pingPacket.Message, time.Unix(pingPacket.SentTime, 0), addr)
		}

		CheckError(err, "Couldn't listen ")
	}
	
}

func (network *Network) SendKademliaPacket(targetNode *Contact, procedure string) {
	//establish a connection to the remote server.
	conn := connect(network.contact.Address, targetNode.Address)

	kademliaPacket := network.CreateKademliaPacket(network.contact.Address, procedure)
	log.Println(kademliaPacket.SourceAddress)

	data, err := proto.Marshal(kademliaPacket)
	CheckError(err, "Couldn't marshal the message")

	buf := []byte(data)

	_, err = conn.Write(buf)
	CheckError(err, "Couldn't write the message")

}

func (network *Network) CreateKademliaPacket(sourceAddress string, procedure string) *KademliaPacket {

	//check that the procedure is one defined by the constants in this file.
	if procedure != PingReq && procedure != PingResp {
		log.Println("bad procedure.." + procedure)
	}

	kademliaPacket := KademliaPacket{
		SourceAddress: sourceAddress,
		Procedure: procedure,
	}
	return &kademliaPacket
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
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