package main

import (
	"demo01/checkError"
	"log"
	"net"
)

func connect2TcpServer(serverAddr string) *net.TCPConn {
	tcpAddre, err := net.ResolveTCPAddr("tcp", serverAddr)
	checkError.CheckError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddre)
	checkError.CheckError(err)
	return conn
}

func SendTcpServer(conn *net.TCPConn) {
	n, err := conn.Write([]byte("hello world"))
	checkError.CheckError(err)
	log.Printf("sending %d bytes", n)
}

func TcpLongConnection() {
	conn := connect2TcpServer("127.0.0.1:5678")
	defer conn.Close()
	for i := 0; i < 10; i++ {
		SendTcpServer(conn)
	}
	log.Println("close connection")
}

func main() {
	TcpLongConnection()
}
