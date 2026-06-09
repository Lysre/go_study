package main

import (
	"bytes"
	"demo01/checkError"
	"demo01/common"
	"io"
	"log"
	"net"
	"time"
)

// log打印表示使用buffer中的内容，重置buffer后，继续接收数据

func TcpStick() {
	TcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:5678")
	checkError.CheckError(err)
	listener, err := net.ListenTCP("tcp4", TcpAddr)
	checkError.CheckError(err)
	log.Printf("waiting for client connection...")
	conn, err := listener.Accept()
	log.Printf("eatablish connection to client %s\n", conn.RemoteAddr().String())
	defer conn.Close()

	time.Sleep(5 * time.Second)
	request := make([]byte, 1024)
	buffer := bytes.Buffer{}

	for {
		n, err := conn.Read(request)
		if err != nil {
			if err == io.EOF {
				if buffer.Len() > 0 {
					log.Println(buffer.String())
				}
			} else {
				checkError.CheckError(err)
			}
			break
		}
		log.Printf("receive request: %s\n", string(request[:n])) // 约定要的分割符号为common.MAGIC
		data := request[:n]
		for { // data 中可能包含多个分隔符号
			pos := bytes.Index(data, common.MAGIC)
			if pos >= 0 {
				if pos == 0 {
					if buffer.Len() > 0 {
						log.Println(buffer.String()) // 把buffer中的内容使用
						buffer.Reset()               // 重置buffer
					}
				} else if pos > 0 {
					buffer.Write(data[:pos])     //将分隔符之前的内容写入buffer中
					log.Println(buffer.String()) //使用buffer中的内容
					buffer.Reset()               // 重置buffer
				}
				data = data[pos+len(common.MAGIC):] // 跳过分隔符
			} else {
				buffer.Write(data) // 在这段内容中没有分隔符，把所有内容写入buffer中
				break
			}
		}
	}
}

func main() {
	TcpStick()
}
