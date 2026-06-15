package main

import (
	"demo01/checkError"
	"demo01/common"
	"encoding/json"
	"log"
	"math/rand/v2"
	"net"
	"sync"
	"time"
)

func connect2UdpServer(addr string) net.Conn {
	conn, err := net.DialTimeout("udp", addr, 5*time.Second)
	checkError.CheckError(err)
	log.Printf("establish connection to server %s\n", conn.RemoteAddr().String())
	return conn
}

func UdpRpcClient() {
	const P = 50
	const C = 10

	wg := sync.WaitGroup{}
	wg.Add(C)
	for i := 0; i < C; i++ {
		go func() {
			defer wg.Done()
			conn := connect2UdpServer("127.0.0.1:5678")
			for j := 0; j < P; j++ {
				request := common.AddRequest{
					RequestId: rand.Int(),
					A:         int(rand.Int()) % 100,
					B:         int(rand.Int()) % 100,
				}

				bs, err := json.Marshal(request)
				if err != nil {
					log.Println("marshal request error:", err)
					continue
				}

				if _, err := conn.Write(bs); err == nil {
					log.Printf("send request, id %d a %d b %d", request.RequestId, request.A, request.B)
				}

			}
			buffer := make([]byte, 1024)
			for j := 0; j < P; j++ {
				n, err := conn.Read(buffer)
				if err != nil {
					log.Println("read error:", err)
					continue
				}
				var response common.AddResponse
				err = json.Unmarshal(buffer[:n], &response)
				if err == nil {
					log.Printf("Sum=%d ResponseId=%d\n", response.Sum, response.ResponseId)
				} else {
					log.Println("unmarshal error:", err)
				}
			}
		}()

	}
	wg.Wait()
}

func main() {

	UdpRpcClient()
}
