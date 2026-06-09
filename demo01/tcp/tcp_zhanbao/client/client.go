package main

import (
	"demo01/checkError"
	"demo01/common"
	"log"
	"net"
)

func Connect2TcpServer() *net.TCPConn {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:5678")
	checkError.CheckError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	return conn
}

func SendTcpServer(conn *net.TCPConn) {
	n, err := conn.Write(append([]byte("hello world"), common.MAGIC...))
	checkError.CheckError(err)
	log.Printf("send %d bytes to server\n", n)
}

func TcpClient() {
	conn := Connect2TcpServer()
	for i := 0; i < 10; i++ {
		SendTcpServer(conn)
	}
	defer conn.Close()
	log.Println("client connected")
}
