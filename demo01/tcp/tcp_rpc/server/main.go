package main

import (
	"bytes"
	myuitl "demo01/tcp/uitl"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

func StartServer(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("服务监听失败:%v", err)
	}
	defer listener.Close()
	fmt.Printf("TCP服务端启动成功,监听地址:%s", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("服务端接受连接失败:%v", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	//设置读写超时
	_ = conn.SetReadDeadline(time.Now().Add(myuitl.ReadWriteTimeout))
	_ = conn.SetWriteDeadline(time.Now().Add(myuitl.ReadWriteTimeout))

	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("新客户端链接:%s", clientAddr)

	// 读取4字节长度头
	lenBuf := make([]byte, 4)
	_, err := io.ReadFull(conn, lenBuf) // 读取4字节长度头
	if err != nil {
		log.Printf("服务端读取长度头失败:%v", err)
		sendErrorResp(conn, "99", "报文解析错误")
		return
	}

	// 读取包体的长度
	var bodyLen uint32
	_ = binary.Read(bytes.NewBuffer(lenBuf), binary.BigEndian, &bodyLen)

	// 读取完整包体
	bodyBuf := make([]byte, bodyLen)
	if _, err := io.ReadFull(conn, bodyBuf); err != nil {
		log.Printf("服务端读取包体失败:%v", err)
		sendErrorResp(conn, "99", "报文解析错误")
		return
	}

	// 拼接完整报文并解包
	fullPkg := append(lenBuf, bodyBuf...) // 拼接完整报文
	header, xmlGbk, _, err := myuitl.UnPackMsg(fullPkg)
	if err != nil {
		log.Printf("解包失败:%v", err)
		sendErrorResp(conn, "BK", "报文解析错误")
		return
	}

	// 打印请求信息
	fmt.Printf("收到请求| 机构编号：%s | 交易码：%s | 流水号:%s\n",
		strings.TrimSpace(string(header.MerCode[:])),
		strings.TrimSpace(string(header.TradeCode[:])),
		strings.TrimSpace(string(header.SerialNo[:])))

	// 解析XML报文
	var reqXml myuitl.ReqXml
	if len(xmlGbk) > 0 {
		xmlUTF8, err := myuitl.GBKToUTF8(xmlGbk)
		if err != nil {
			log.Printf("GBK解码xml报文失败:%v", err)
			sendErrorResp(conn, "BK", "报文解析错误")
			return
		}
		if err = xml.Unmarshal([]byte(xmlUTF8), &reqXml); err != nil {
			log.Printf("XML解析xml报文失败:%v", err)
			sendErrorResp(conn, "BK", "报文解析错误")
			return
		}
		fmt.Printf("解析xml报文成功:%+v\n", reqXml)
	}

	// 模拟业务处理
	respHeader := myuitl.MsgHeader{}
	myuitl.FillHeader(&respHeader,
		"0000",
		strings.TrimSpace(string(header.MerCode[:])),   // 机构号
		strings.TrimSpace(string(header.TradeCode[:])), // 交易码
		strings.TrimSpace(string(header.SerialNo[:])),  // 流水号
		"s",  // 消息标志
		"00", // 错误码
		"")   // 错误信息
	respXml := myuitl.RespXml{
		Result: "00",
		Data:   "success",
	}

	// 打包并发送应答
	respPkg, err := myuitl.PackMsg(respHeader, respXml, true)
	if err != nil {
		log.Printf("向客户端[%s]打包应答失败:%v", clientAddr, err)
		sendErrorResp(conn, "99", "应答打包失败")
		return
	}

	if _, err = conn.Write(respPkg); err != nil {
		log.Printf("向客户端[%s]发送应答失败:%v", clientAddr, err)
		return
	}

	fmt.Printf("客户端[%s]应发发送成功，链接关闭\n", clientAddr)
}

// 发送异常应答（仅报文头 无 xml 报文 无签名字段）
func sendErrorResp(conn net.Conn, errCode, errMsg string) {
	var respHeader myuitl.MsgHeader
	myuitl.FillHeader(&respHeader, "0000", "", "", "", "s", errCode, errMsg)
	respPkg, _ := myuitl.PackMsg(respHeader, nil, false)
	_, _ = conn.Write(respPkg)
}

func main() {

}
