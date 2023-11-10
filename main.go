package main

import (
	"fmt"
	"os"
)

func main() {
	var dirPath string
	if len(os.Args) > 1 && os.Args[1] != "" {
		dirPath = os.Args[1]
	} else {
		dirPath = "."
	}

	if dirPath == "" {
		dirPath = "."
	}

	filePath := dirPath + "/go.mod"

	data, err := os.ReadFile(filePath)
	if err == nil {
		if err = Run(data, []string{"no-replace"}); err == nil {
			fmt.Print("go.mod is OK\n")
			return
		}
	}

	var returnCode int
	if errWithCode, ok := err.(interface{ Code() int }); ok {
		returnCode = errWithCode.Code()
	} else {
		returnCode = 1
	}
	fmt.Println("go.mod is invalid:", err)
	os.Exit(returnCode)
}
