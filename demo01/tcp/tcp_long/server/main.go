package main

import (
	"demo01/checkError"
	"log"
	"net"
)

func TcpLongConnection() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:5678")
	checkError.CheckError(err)
	listen, err := net.ListenTCP("tcp4", tcpAddr)
	checkError.CheckError(err)
	log.Printf("waiting for client connection...")
	conn, err := listen.Accept()
	checkError.CheckError(err)
	log.Printf("eatablish connection to client %s\n", conn.RemoteAddr().String())
	//conn.SetDeadline(time.Now().Add(10 * time.Second))
	defer conn.Close()

	request := make([]byte, 1024)
	for {
		n, err := conn.Read(request)
		checkError.CheckError(err)
		log.Printf("receive %s", string(request[:n]))
		//conn.SetDeadline(time.Now().Add(10 * time.Second)) // 重置deadline
	}
}

func main() {
	TcpLongConnection()
}
