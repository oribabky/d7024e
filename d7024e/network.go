package d7024e
import (
	"time"
	"log"
	"github.com/golang/protobuf/proto"
	//"fmt"
	)

type Network struct {
}


func (network *Network) Listen(ip string, port int) {
	buf := make([]byte, 1024)
	serverConn := listening(ip)
	for {
		n, addr, err := serverConn.ReadFromUDP(buf)
		pingPacket := &PingPacket{}
		err = proto.Unmarshal(buf[0:n], pingPacket)
		log.Printf("Received %s at %s from %s", *pingPacket.Message, time.Unix(*pingPacket.SentTime, 0), addr)
		CheckError(err, "Couldn't listen ")
	}
	
}

func (network *Network) SendPingMessage(sender *Contact, remote *Contact) {
	//establish a connection to the remote server.
	conn := connect(sender.Address, remote.Address)

	pingPacket := network.CreatePingPacket("Hello I'm alive")
	now := time.Now().Unix()
	pingPacket.SentTime = &now

	data, err := proto.Marshal(pingPacket)
	CheckError(err, "Couldn't marshal the message")


	_, err = conn.Write(data)
	CheckError(err, "Couldn't write the message")

}

func (network *Network) CreatePingPacket(msg string) *PingPacket {
	pingPacket := PingPacket{
		Message: &msg,
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