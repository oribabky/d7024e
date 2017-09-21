package d7024e
import (
	"time"
	"log"
	"github.com/golang/protobuf/proto"
	//"fmt"
	)

type Network struct {
	contact *Contact
}

func NewNetwork(contact *Contact) Network {
	return Network{contact}
}

func (network *Network) Listen() {
	buf := make([]byte, 1024)
	serverConn := listening(network.contact.Address)
	log.Println(network.contact.Address)
	for {
		n, addr, err := serverConn.ReadFromUDP(buf)
		pingPacket := &PingPacket{}
		err = proto.Unmarshal(buf[0:n], pingPacket)
		if addr != nil {
			log.Printf("Received %s at %s from %s", pingPacket.Message, time.Unix(pingPacket.SentTime, 0), addr)
		}

		CheckError(err, "Couldn't listen ")
	}
	
}

func (network *Network) SendPingMessage(remote *Contact) {
	//establish a connection to the remote server.
	conn := connect(network.contact.Address, remote.Address)

	pingPacket := network.CreatePingPacket(network.contact.Address)

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