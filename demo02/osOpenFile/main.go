package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.OpenFile("../data/test.txt", os.O_CREATE|os.O_RDONLY|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Open file failed:%s\n", err)
	} else {
		defer file.Close()
		file.WriteString("hello world\n")
		file.WriteString("Lysre\n")
	}
}
