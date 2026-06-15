package main

import (
	"demo01/checkError"
	"net"
)

func TcpServer(address string) net.Conn {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	checkError.CheckError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	conn, err := listener.Accept()
	checkError.CheckError(err)
	return conn
}

func ReadTcpRsa(conn net.Conn) {
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	checkError.CheckError(err)

}

func TcpRsa() {
	conn := TcpServer("127.0.0.1:8443")
	defer conn.Close()
	ReadTcpRsa(conn)
}

func main() {
	TcpRsa()
}
