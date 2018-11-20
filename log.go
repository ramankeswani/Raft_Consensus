package main

import (
	"fmt"
	"os"
	"strings"
)

var f *os.File
var t string

func initLog(nodeID string) {
	file, err := os.Create("logs/" + nodeID + ".txt")
	checkError(err, "initLog")
	f = file
	t = "recover"
}

func logFile(tag string, message string) {
	if !strings.Contains(message, "ThisIsHeartbeat") && strings.Compare(tag, t) == 0 {
		fmt.Fprintf(f, message)
	}
}
