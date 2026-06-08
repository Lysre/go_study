package checkError

import (
	"fmt"
	"os"
)

func CheckError(err error) {
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
}
