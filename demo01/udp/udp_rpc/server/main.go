package main

import (
	"demo01/checkError"
	"demo01/common"
	"encoding/json"
	"io"
	"log"
	"net"
)

func UdpRpcServer() {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:5678")
	checkError.CheckError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	checkError.CheckError(err)
	log.Println("start udp server")
	defer conn.Close()

	//===================接收客户端请求===================
	const P = 1000
	for i := 0; i < P; i++ {
		go func() {
			for {
				request := make([]byte, 1024)
				n, remoteAddr, err := conn.ReadFromUDP(request)
				if err != nil && err != io.EOF {
					log.Println("read error:", err)
					continue
				}
				response := HandleRequest(request[:n])
				if len(response) > 0 {
					conn.WriteToUDP(response, remoteAddr)
				}
			}
		}()
	}
	select {}
}

func HandleRequest(request []byte) (response []byte) {
	var requestMsg common.AddRequest
	err := json.Unmarshal(request, &requestMsg)
	checkError.CheckError(err)
	log.Printf("requestId=%d a=%d b=%d\n", requestMsg.RequestId, requestMsg.A, requestMsg.B)
	responseMsg := common.AddResponse{
		ResponseId: requestMsg.RequestId,
		Sum:        requestMsg.A + requestMsg.B,
	}
	response, err = json.Marshal(responseMsg)
	checkError.CheckError(err)
	return
}

func main() {
	UdpRpcServer()
}
