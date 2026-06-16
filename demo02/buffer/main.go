package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// ReadFileWithBuffer 读取文件，使用缓冲区
func ReadFileWithBuffer(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Open file failed:%s\n", err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n') // 以字符串的方式读取，参数为分隔符
		if len(line) > 0 {
			fmt.Print(line)
		}
		if err == io.EOF {
			break
		}
	}
}

// WriteFileWithBuffer 带缓冲区的写入文件
func WriteFileWithBuffer(filePath string) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("Open file failed:%s\n", err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	writer.Write([]byte("hello world\n"))
	writer.WriteString("Lysre\n")
	writer.Flush() // 刷新缓冲区，将缓冲区中的数据写入文件，应为缓冲区默认大小为1024，所以需要刷新缓冲区才能将数据写入文件
}

func main() {
	WriteFileWithBuffer("../data/test.txt")
	ReadFileWithBuffer("../data/test.txt")
}
