package main

import (
	"demo01/checkError"
	"log"
	"net"
)

func connect2TcpServer(serverAddr string) *net.TCPConn {
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
	checkError.CheckError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError.CheckError(err)
	return conn
}

func SendTcpServer(conn *net.TCPConn) {
	n, err := conn.Write([]byte("hello world"))
	checkError.CheckError(err)
	log.Printf("send %d bytes\n", n)
}

func TcpClient() {
	conn := connect2TcpServer("127.0.0.1:5678")
	defer conn.Close()
	SendTcpServer(conn)
}

func main() {
	TcpClient()
}
