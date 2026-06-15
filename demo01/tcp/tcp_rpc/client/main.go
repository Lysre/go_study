package main

import (
	"bytes"
	myuitl "demo01/tcp/uitl"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

func SendTCPRequest(addr string, merCode, tradeCode string, reqXml myuitl.ReqXml) error {
	//建立短链接
	conn, err := net.DialTimeout("tcp", addr, time.Second*3)
	if err != nil {
		return fmt.Errorf("链接服务端失败: %w", err)
	}
	defer conn.Close()
	_ = conn.SetReadDeadline(time.Now().Add(myuitl.ReadWriteTimeout))
	_ = conn.SetWriteDeadline(time.Now().Add(myuitl.ReadWriteTimeout))

	// 组装报文头
	var reqHeader myuitl.MsgHeader
	serialNo := myuitl.GenerateSerialNo(0) // 生成流水号
	myuitl.FillHeader(&reqHeader,
		"0000",
		merCode,
		tradeCode,
		serialNo,
		"q",
		" ", //请求错误码填空格
		" ") //请求错误信息填空格

	// 打包请求报文
	sendPkg, err := myuitl.PackMsg(reqHeader, reqXml, true)
	if err != nil {
		return fmt.Errorf("打包请求报文失败:%w", err)
	}

	// 发送请求
	_, err = conn.Write(sendPkg)
	if err != nil {
		return fmt.Errorf("发送请求报文失败:%w", err)
	}
	fmt.Printf("请求发送成功 | 流水号:%s\n", serialNo)
	fmt.Println("等待服务端应答。。。")

	// 读取应答长度
	lenBuf := make([]byte, 4)
	_, err = io.ReadFull(conn, lenBuf)
	if err != nil {
		return fmt.Errorf("读取应答长度失败:%w", err)
	}

	var bodyLen uint32
	_ = binary.Read(bytes.NewBuffer(lenBuf), binary.LittleEndian, &bodyLen) // 读取应答长度

	// 读取应答包体
	bodyBuf := make([]byte, bodyLen)
	_, err = io.ReadFull(conn, bodyBuf)
	if err != nil {
		return fmt.Errorf("读取应答包体失败:%w", err)
	}

	// 解包应答
	respFull := append(lenBuf, bodyBuf...)
	respHeader, xmlGBK, _, err := myuitl.UnPackMsg(respFull)
	if err != nil {
		return fmt.Errorf("解包应答报文失败:%w", err)
	}

	// 处理应答报文头
	errCode := strings.TrimSpace(string(respHeader.ErrCode[:]))
	errMsg := strings.TrimSpace(string(respHeader.ErrMsg[:]))

	if errCode != "00" {
		return fmt.Errorf("交易失败:%s", errMsg)
	}

	// 解析应答xml
	if len(xmlGBK) > 0 {
		xmlUtf8, err := myuitl.GBKToUTF8(xmlGBK)
		if err != nil {
			return fmt.Errorf("解析应答xml失败:%w", err)
		}
		var respXml myuitl.ReqXml
		err = xml.Unmarshal([]byte(xmlUtf8), &respXml)
		if err != nil {
			return fmt.Errorf("解析应答xml失败:%w", err)
		}
		fmt.Printf("应答xml:%v\n", respXml)
	}
	return nil
}
