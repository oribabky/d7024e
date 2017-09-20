package d7024e

import {
	"flag"
	"log"
	"net"
}

connect(senderAddr string, destinationAddr string) {
	remoteAddr, err := net.ResolveUDPAddr("udp", destinationAddr)
	CheckError(err, "")

	localAddr, err := net.ResolveUDPAddr("udp", senderAddr)
	CheckError(err, "")

	conn, err := net.DialUDP("udp", localAddr, remoteAddr)
	CheckError(err, "")

	//if there is an error, close the connection
	defer conn.Close()
}

listening(localAddr string) {
	serverAddr, err := net.ResolveUDPAddr("udp", localAddr)
	CheckError(err, "")

	serverConn, err := net.ListenUDP("udp", serverAddr)
	CheckError(err, "")
	defer serverConn.Close()
}