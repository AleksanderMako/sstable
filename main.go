package main

import (
	"fmt"
	"os"
)

func main() {

	logsDir, err := os.Getwd()
	if err != nil {
		// do smt
	}
	os.Setenv("LOGS_DIR", logsDir+"/logs")
	fmt.Println(os.Getenv("LOGS_DIR"))
}
