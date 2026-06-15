package myuitl

import (
	"bytes"
	"crypto/dsa"
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	HeaderTotalLen    = 148              // 报文头总字节数
	SignFieldTotalLen = 3 + 96 + 33      // 签名字段总字节数 前缀+签名+后缀
	ReadWriteTimeout  = 60 * time.Second // 读写超时时间
	PadChar           = ' '              // 填充字符
)

type MsgHeader struct {
	Version   [4]byte   //4 版本号 0000
	MerCode   [15]byte  //15机构号
	TradeCode [10]byte  // 10交易码
	SerialNo  [16]byte  // 16流水号
	MsgFlag   [1]byte   // 1消息标志 q=请求 s=响应
	ErrCode   [2]byte   // 2错误码
	ErrMsg    [100]byte // 100 错误信息
}

// FillHeader 填充报文头，不足的地方补空格
func FillHeader(h *MsgHeader, version, merCode, tradeCode, serialNo, msgFlag, errCode, errMsg string) {
	copy(h.Version[:], PadString(version, 4))
	copy(h.MerCode[:], PadString(merCode, 15))
	copy(h.TradeCode[:], PadString(tradeCode, 10))
	copy(h.SerialNo[:], PadString(serialNo, 16))
	copy(h.MsgFlag[:], PadString(msgFlag, 1))
	copy(h.ErrCode[:], PadString(errCode, 2))
	copy(h.ErrMsg[:], PadString(errMsg, 100))
}

// GenerateSerialNo 生成标准流水号：YYYYMMDD + 8位序号
func GenerateSerialNo(seq int) string {
	date := time.Now().Format("20060102")
	return fmt.Sprintf("%s%08d", date, seq)
}

// PadString 字符串定长不足补空格，超出定长截取
func PadString(s string, length int) []byte {
	b := make([]byte, length)
	copy(b, []byte(s))
	for i := len(s); i < length; i++ {
		b[i] = PadChar
	}
	return b
}

// UTF8ToGBK GBK 编解码
// @param s	待编码的字符串
func UTF8ToGBK(s string) ([]byte, error) { // UTF-8 编码为 GBK
	reader := transform.NewReader(bytes.NewReader([]byte(s)), simplifiedchinese.GBK.NewEncoder())
	return io.ReadAll(reader)
}

// GBKToUTF8 GBK 编解码
// @param b	待解码的GBK字节数组
func GBKToUTF8(b []byte) (string, error) { // GBK 解码为 UTF-8
	reader := transform.NewReader(bytes.NewBuffer(b), simplifiedchinese.GBK.NewDecoder())
	d, err := io.ReadAll(reader)
	return string(d), err
}

// ReqXml 请求Xml报文结构体
type ReqXml struct {
	XMLName xml.Name `xml:"Req"`
	MerCode string   `xml:"mer_code"`
	TradeNo string   `xml:"trade_no"`
	Amount  string   `xml:"amount"`
}

// RespXml 响应Xml报文结构体
type RespXml struct {
	XMLName xml.Name `xml:"Resp"`
	Result  string   `xml:"result"`
	Data    string   `xml:"data,omitempty"`
}

// 签名计算	MD5 + DSA

var (
	DSAPrivateKey *dsa.PrivateKey // DSA 私钥
	DSAPublicKey  *dsa.PublicKey  // DSA 公钥
)

// InitDsaKey 初始化DSA密钥对
func InitDsaKey() error { // 初始化 DSA 私钥
	params := new(dsa.Parameters)                                     // DSA 参数
	err := dsa.GenerateParameters(params, rand.Reader, dsa.L1024N160) // 生成 DSA 参数
	if err != nil {
		return err
	}

	priv := new(dsa.PrivateKey) // DSA 私钥
	priv.Parameters = *params
	err = dsa.GenerateKey(priv, rand.Reader)
	if err != nil {
		return nil
	}
	DSAPrivateKey = priv
	DSAPublicKey = &priv.PublicKey
	return nil
}

// CalcMD5 计算MD5(报文头二进制 + GBK-XML) 返回16进制字符串
// @param headBin	报文头二进制字节数组
// @param xmlGbk	GBK-XML二进制字节数组
func CalcMD5(headBin, xmlGbk []byte) string {
	h := md5.New()
	h.Write(headBin)                      // 写入报文头
	h.Write(xmlGbk)                       // 写入 xml 报文
	return hex.EncodeToString(h.Sum(nil)) // 返回 MD5 值的十六进制字符串
}

// DsaSign DSA签名，输入原始数据，返回16进制签名字符串
// @param md5Str	MD5 值的十六进制字符串
func DsaSign(md5Str string) (string, error) {
	r, s, err := dsa.Sign(rand.Reader, DSAPrivateKey, []byte(md5Str))
	// 得到两个独立大数 r、s，都是 *big.Int
	if err != nil {
		return "", err
	}

	//	拼接 r,s
	// r.Bytes()：r 转二进制字节
	// s.Bytes()：s 转二进制字节
	// 拼接成一段连续二进制流
	sigBin := append(r.Bytes(), s.Bytes()...)
	return hex.EncodeToString(sigBin), nil
}

// BuildSignField 组装签名字段：3位长度 + 签名串 + 补空格至132字节长度
// @param sigStr DSA签名字符串
func BuildSignField(sigStr string) ([]byte, error) {
	sigLen := len(sigStr)
	if sigLen < 92 || sigLen > 96 {
		return nil, errors.New("DSA签名长度非法，要求92~96字节")
	}

	lenPrefix := fmt.Sprintf("%03d", sigLen) // 3位长度前缀 不足左补0
	// 拼接前缀+签名
	full := lenPrefix + sigStr
	// 补空格至132字节长度
	padCnt := SignFieldTotalLen - len(full)
	if padCnt < 0 {
		return nil, errors.New("签名字段过长")
	}
	full += strings.Repeat(string(PadChar), padCnt) // strings.Repeat()：重复字符串 padCnt 次
	return []byte(full), nil                        // full转换为字节数组，二进制编码
}

// 报文打包/解包

// PackMsg 打包完整报文
// @param header 报文头
// @param xmlobj xml 报文结构体
// @param hasXml 是否有 xml 报文
func PackMsg(header MsgHeader, xmlObj any, hasXml bool) ([]byte, error) {
	var xmlGbk []byte
	var err error

	// 1.序列化报文头为二进制
	headBuf := new(bytes.Buffer)                          // new 创建一个 bytes.Buffer 实例
	err = binary.Write(headBuf, binary.BigEndian, header) // 写入报文头
	if err != nil {
		return nil, fmt.Errorf("序列化报文头失败:%w", err)
	}

	headBin := headBuf.Bytes() // 报文头二进制
	if len(headBin) != HeaderTotalLen {
		return nil, fmt.Errorf("报文头长度异常，实际长度:%d，要求长度:%d", len(headBin), HeaderTotalLen)
	}

	// 处理xml报文
	if hasXml {
		xmlUtf8, err := xml.Marshal(xmlObj)
		if err != nil {
			return nil, fmt.Errorf("序列化xml报文失败:%w", err)
		}
		xmlGbk, err = UTF8ToGBK(string(xmlUtf8))
		if err != nil {
			return nil, fmt.Errorf("GBK编码xml报文失败:%w", err)
		}
	}

	// 生成签名字段
	var signField []byte
	if hasXml {
		// 正常报文：计算MD5->DSA签名->组装签名字段
		md5Str := CalcMD5(headBin, xmlGbk)
		dsaSigStr, err := DsaSign(md5Str)
		if err != nil {
			return nil, fmt.Errorf("DSA签名失败:%w", err)
		}
		signField, err = BuildSignField(dsaSigStr)
		if err != nil {
			return nil, fmt.Errorf("组装签名字段失败:%w", err)
		}
	} else {
		// 异常报文：签名字段全补空格
		signField = make([]byte, SignFieldTotalLen)
		for i := range signField {
			signField[i] = PadChar
		}
	}

	// 计算报文总长度：head+xml+signField(不包含靠头4字节)
	bodyLen := HeaderTotalLen + len(xmlGbk) + SignFieldTotalLen

	// 拼接完整报文
	var pkgBuf bytes.Buffer
	_ = binary.Write(&pkgBuf, binary.BigEndian, &bodyLen)
	pkgBuf.Write(headBin)
	pkgBuf.Write(xmlGbk)
	pkgBuf.Write(signField)
	return pkgBuf.Bytes(), nil
}

// UnPackMsg 解包完成报文
// @param raw 完整报文二进制字节数组
func UnPackMsg(raw []byte) (header MsgHeader, xmlGbk []byte, signField []byte, err error) {
	buf := bytes.NewBuffer(raw)

	// 读取4个字节总长度
	var bodyLen uint32
	err = binary.Read(buf, binary.BigEndian, &bodyLen)
	if err != nil {
		err = fmt.Errorf("读取总长度失败:%w", err)
		return
	}

	bodyData := make([]byte, bodyLen)
	if _, err := io.ReadFull(buf, bodyData); err != nil {
		err = fmt.Errorf("读取报文体失败:%w", err)
		return
	}

	if len(bodyData) < HeaderTotalLen+SignFieldTotalLen {
		err = errors.New("报文体长度不足")
		return
	}

	headBin := bodyData[:HeaderTotalLen]
	signField = bodyData[len(bodyData)-SignFieldTotalLen:]
	xmlGbk = bodyData[HeaderTotalLen : len(bodyData)-SignFieldTotalLen]

	// 反序列化报文头
	headBuf := bytes.NewBuffer(headBin)
	err = binary.Read(headBuf, binary.BigEndian, &header)
	if err != nil {
		err = fmt.Errorf("反序列化报文头失败:%w", err)
		return
	}
	return header, xmlGbk, signField, nil
}
