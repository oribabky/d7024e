package d7024e
import (
	//"time"
	"log"
	"github.com/golang/protobuf/proto"
	//"fmt"
	"net"
	)

type Network struct {
	contact *Contact
}

func NewNetwork(contact *Contact) Network {
	return Network{contact}
}

func (network *Network) Listen() {
	buf := make([]byte, 1024)

	//establish a connection 
	serverAddr, err := net.ResolveUDPAddr("udp", network.contact.Address)
	CheckError(err, "")
	serverConn, err := net.ListenUDP("udp", serverAddr)
	CheckError(err, "")
	defer serverConn.Close() //close the connection when something is return

	//log.Println(network.contact.Address)
	for {
		n, addr, err := serverConn.ReadFromUDP(buf)
		pingPacket := &PingPacket{}
		err = proto.Unmarshal(buf[0:n], pingPacket)
		if addr != nil {
			log.Printf("Received ping from %s", pingPacket.Message) // must be changed !!!!!
		}

		CheckError(err, "Couldn't listen ")
	}
	
}

func (network *Network) SendPingMessage(remote *Contact) {
	//establish a connection to the remote server.
	remoteAddr, err := net.ResolveUDPAddr("udp", remote.Address)
	CheckError(err, "")
	localAddr, err := net.ResolveUDPAddr("udp", network.contact.Address)
	CheckError(err, "")
	conn, err := net.DialUDP("udp", localAddr, remoteAddr)
	CheckError(err, "")
	defer conn.Close() //if there is an error, close the connection
	//conn := connect(network.contact.Address, remote.Address)

	pingPacket := network.CreatePingPacket(network.contact.Address)  //here we can set the ping message

	//now := time.Now().Unix()
	//pingPacket.SentTime = now

	data, err := proto.Marshal(pingPacket)
	CheckError(err, "Couldn't marshal the message")

	buf := []byte(data)

	_, err = conn.Write(buf)
	CheckError(err, "Couldn't write the message")

}

func (network *Network) CreatePingPacket(msg string) *PingPacket {
	pingPacket := PingPacket{
		Message: msg,
	}
	return &pingPacket
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