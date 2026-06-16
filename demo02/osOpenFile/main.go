package main

import (
	"fmt"
	"io"
	"os"
)

// WriteFile os.O_CREATE 创建文件
// os.O_RDONLY 只读
// os.O_WRONLY 只写
// os.O_APPEND 追加写入
// os.O_SYNC 同步写入
// os.O_EXCL 排他写入
// os.O_SYNC 同步写入
// os.O_TRUNC 截断文件
func WriteFile(filePath string) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDONLY|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Open file failed:%s\n", err)
	} else {
		defer file.Close()
		file.WriteString("hello world\n")
		file.WriteString("Lysre\n")
	}
}

func ReadFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Open file failed:%s\n", err)
	}
	defer file.Close()
	bs := make([]byte, 10) // (类型，长度，初始值)	长度为0时是读不出信息的
	file.Read(bs)
	fmt.Println(string(bs))

	file.Seek(10, 1) //  offset 表示偏移量，whence 表示偏移量的参考点; 0 表示从文件开头读取， 1 表示从文件当前位置读取， 2 表示从文件末尾读取

}

func ReadBigFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Open file failed:%s\n", err)
	}
	defer file.Close()
	bs := make([]byte, 10)
	for {
		n, err := file.Read(bs)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Read file failed:%s\n", err)
			return
		}
		if n > 0 {
			fmt.Println(string(bs[:n]))
		}

	}
}

func main() {
	WriteFile("../data/test.txt")
	ReadBigFile("../data/test.txt")
}
