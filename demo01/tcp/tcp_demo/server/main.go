package main

import (
	"demo01/checkError"
	"log"
	"net"
	"time"
)

/*
   我们需要建立链接 TCP 是一对一的链接所以是需要阻塞的
*/

func TCPServer() {
	addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:5678") //讲地址转换成 结构体
	checkError.CheckError(err)                                // 自定义错误

	listen, err := net.ListenTCP("tcp4", addr) // 建立 tcp 链接
	checkError.CheckError(err)
	log.Println("waiting for client connection...")

	conn, err := listen.Accept() // 对链接阻塞等待消息
	checkError.CheckError(err)

	log.Printf("eatablish connection to client %s\n", conn.RemoteAddr().String())
	conn.SetReadDeadline(time.Now().Add(3 * time.Second)) //设置一个读的期限，超过这个期限再调Read()就会发生error。默认是60s内可Read()。
	defer conn.Close()

	request := make([]byte, 1024)
	n, err := conn.Read(request)
	log.Printf("receive %s\n", string(request[:n]))
}

func main() {
	TCPServer()
}
