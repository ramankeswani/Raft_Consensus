package main

import (
	"fmt"
	"os"
	"strings"
)

var f *os.File

func initLog(nodeID string) {
	file, err := os.Create("logs/" + nodeID + ".txt")
	checkError(err, "initLog")
	f = file
}

func logFile(tag string, message string) {
	if strings.Compare(tag, logTag) == 0 && !strings.Contains(message, "ThisIsHeartbeat") {
		fmt.Fprintf(f, message)
	}
}
