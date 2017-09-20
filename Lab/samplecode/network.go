package d7024e

import (
	"time"
	"log"
	"github.com/golang/protobuf/proto"
	)

type Network struct {
}


func Listen(ip string, port int) {
	listening(ip)
	//n, addr, err 
	// TODO
}

func (network *Network) SendPingMessage(sender *Contact, remote *Contact) {
	//establish a connection to the remote server.
	connect(sender.Address, remote.Address)

	pingPacket := CreatePingPacket("Hello I'm alive")
	pingPacket.sent_time := time.Now().Unix()

	data, err := proto.Marshal(pingPacket)
	CheckError(err, "Couldn't marshal the message")


	_, err = conn.Write(data)
	CheckError(err, "Couldn't write the message")

}

func (network *Network) CreatePingPacket(msg string) *PingPacket {
	pingPacket := PingPacket{
		message: &msg,
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