package main

import (
	"demo01/checkError"
	"log"
	"net"
	"sync"
)

func UdpServer() {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:5678")
	checkError.CheckError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	checkError.CheckError(err)
	defer conn.Close()

	request := make([]byte, 1024)
	for i := 0; i < 10; i++ {
		n, remoteAddr, err := conn.ReadFromUDP(request)
		checkError.CheckError(err)
		log.Printf("receive %s from %s\n", string(request[:n]), remoteAddr.String())
	}
}

// UdpServerCurrent 服务器端并发处理多个客户端的请求
func UdpServerCurrent() {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:5678")
	checkError.CheckError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	checkError.CheckError(err)
	defer conn.Close()

	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			request := make([]byte, 1024)
			for i := 0; i < 10; i++ {
				n, remoteAddr, err := conn.ReadFromUDP(request)
				checkError.CheckError(err)
				log.Printf("syncId %d receive %s from %s\n", i, string(request[:n]), remoteAddr.String()) // 模拟对收到数据的处理
			}
		}(i)
	}
	wg.Wait()

}

func main() {
	UdpServerCurrent()
}
