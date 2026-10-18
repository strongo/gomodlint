package main

import (
	"fmt"
	"os"
)

var osExit = os.Exit

func main() {
	var dirPath string
	if len(os.Args) > 1 {
		dirPath = os.Args[1]
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
	osExit(returnCode)
}
