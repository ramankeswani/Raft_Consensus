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
	t = "*"
}

func logFile(tag string, message string) {
	if strings.Compare(t, "*") == 0 && !strings.Contains(message, "ThisIsHeartbeat") {
		fmt.Fprintf(f, message)
	} else if !strings.Contains(message, "ThisIsHeartbeat") && strings.Compare(tag, t) == 0 {
		fmt.Fprintf(f, message)
	}
}
