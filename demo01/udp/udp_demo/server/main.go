package main

import (
	"demo01/checkError"
	"log"
	"net"
)

// UdpServer 服务器端
// udp链接和tcp不同他一个服务可以连接多个客户端，tcp需要accept()方法接受客户端的连接
// 但是可以通过ReadFromUDP()方法读取客户端发送的消息	会返回一个带有客户端地址的结构体remoteAddr
func UdpServer() {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:5678")
	checkError.CheckError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	checkError.CheckError(err)
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(buffer)
	checkError.CheckError(err)
	log.Printf("receive %s from %s\n", string(buffer[:n]), remoteAddr.String())
}

func main() {
	UdpServer()
}
