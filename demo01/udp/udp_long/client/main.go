package main

import (
	"demo01/checkError"
	"log"
	"net"
	"time"
)

func connect2UdpServer(serverAddr string) net.Conn {
	conn, err := net.DialTimeout("udp", serverAddr, time.Second)
	checkError.CheckError(err)
	log.Printf("establish connection to server %s\n", conn.RemoteAddr().String())
	return conn
}

func SendUdpServer(conn net.Conn) {
	n, err := conn.Write([]byte("hello world"))
	checkError.CheckError(err)
	log.Printf("send %d bytes to server %s\n", n, conn.RemoteAddr().String())
}

func UdpLongConnection() {
	conn := connect2UdpServer("127.0.0.1:5678")
	for i := 0; i <= 3; i++ {
		SendUdpServer(conn)
	}
	conn.Close()
	log.Println("server connection closed")
}

func main() {
	UdpLongConnection()
}
