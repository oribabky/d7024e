package d7024e

import (
	//"flag"
	//"log"
	"net"
)

func connect(senderAddr string, destinationAddr string) {
	remoteAddr, err := net.ResolveUDPAddr("udp", destinationAddr)
	CheckError(err, "")

	localAddr, err := net.ResolveUDPAddr("udp", senderAddr)
	CheckError(err, "")

	conn, err := net.DialUDP("udp", localAddr, remoteAddr)
	CheckError(err, "")

	//if there is an error, close the connection
	defer conn.Close()
}

func listening(localAddr string) *net.UDPConn {
	serverAddr, err := net.ResolveUDPAddr("udp", localAddr)
	CheckError(err, "")

	serverConn, err := net.ListenUDP("udp", serverAddr)
	CheckError(err, "")
	defer serverConn.Close()
	return serverConn
}