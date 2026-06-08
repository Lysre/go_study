package main

import (
	"demo01/checkError"
	"log"
	"net"
	"time"
)

func connect2UdpServer(remoteAddr string) net.Conn {
	conn, err := net.DialTimeout("udp", remoteAddr, time.Second)
	checkError.CheckError(err)
	log.Printf("establish connection to server %s\n", conn.RemoteAddr().String())
	return conn
}

func SendUdpServer(conn net.Conn) {
	n, err := conn.Write([]byte("hello world"))
	checkError.CheckError(err)
	log.Printf("send %d bytes to server\n", n)
}

func UdpClient() {
	conn := connect2UdpServer("127.0.0.1:5678")
	SendUdpServer(conn)
	defer conn.Close()
}

func main() {
	UdpClient()
}
