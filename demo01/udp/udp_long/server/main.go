package main

import (
	"demo01/checkError"
	"log"
	"net"
	"time"
)

func UdpLongConnection() {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:5678")
	checkError.CheckError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	checkError.CheckError(err)
	defer conn.Close()

	time.Sleep(5 * time.Second)
	request := make([]byte, 1024)
	for {
		conn.SetDeadline(time.Now().Add(5 * time.Second))
		n, remoterAddr, err := conn.ReadFromUDP(request)
		checkError.CheckError(err)
		log.Printf("read %d bytes from %s", n, remoterAddr.String())
		log.Println(string(request[:n]))
	}
}

func main() {
	UdpLongConnection()
}
