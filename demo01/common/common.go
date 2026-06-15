package common

import "os"

var MAGIC = []byte{1, 1, 5, 2, 0}

var (
	publicKey  []byte
	privateKey []byte
)

type AddRequest struct {
	RequestId int `json:"request_id"`
	A         int `json:"a"`
	B         int `json:"b"`
}

type AddResponse struct {
	ResponseId int `json:"response_id"`
	Sum        int `json:"sum"`
}

func ReadFile(keyFile string) ([]byte, error) {
	f, err := os.Open(keyFile)
	if err != nil {
		return nil, err
	}
	content := make([]byte, 4096)
	n, err := f.Read(content)
	if err != nil {
		return nil, err
	}
	return content[:n], nil
}
