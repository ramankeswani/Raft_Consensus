package main

import (
	"net"
	"strconv"
	"strings"
)

func initRecovery(port int) {
	logFile("recover", "initRecoveryStarts\n")
	message := nodeID + " " + NodeRecoverMessage + " " + strconv.Itoa(port) + "\n"
	for node := range otherNodes {
		chanMap[otherNodes[node].nodeID] <- message
	}
	logFile("recover", "initRecoveryEnds\n")
}

func handleRecoveryMessage(message string) {
	logFile("recover", "handleRecoveryMessage Starts\n")
	dataSlice := strings.Split(strings.TrimRight(message, "\n"), " ")
	port, _ := strconv.Atoi(dataSlice[2])
	address := "127.0.0.1:" + strconv.Itoa(port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	checkError(err, "client")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	c := connection{
		nodeID: dataSlice[0],
		conn:   conn,
	}
	connMap[dataSlice[0]] = c
	totalNodes++
	logFile("recover", "handleRecoveryMessage Ends\n")
}
