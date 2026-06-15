package main

import (
	"demo01/checkError"
	"log"
	"net"
	"sync"
	"time"
)

func connect2UdpServer(remoteAddr string) net.Conn {
	conn, err := net.DialTimeout("udp", remoteAddr, 3*time.Second)
	checkError.CheckError(err)
	log.Printf("establish connection to server %s\n", conn.RemoteAddr().String())
	return conn
}

func sendUdpServer(conn net.Conn) {
	n, err := conn.Write([]byte("hello world"))
	checkError.CheckError(err)
	log.Printf("send %d bytes to server\n", n)
}

func UdpConnectionCurrent() {
	conn := connect2UdpServer("127.0.0.1:5678")
	defer conn.Close()
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			sendUdpServer(conn)
		}(i)
	}
	wg.Wait()
}

func main() {
	UdpConnectionCurrent()
}
