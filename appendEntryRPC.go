package main

import (
	"strconv"
)

/*
Creates Append Entry RPC Request and sends the request to all Clients
*/
func appendEntryInit(command string) {
	logFile("append", "Append Entry Starts Map len:"+strconv.Itoa(len(connMap))+"\n")
	prevLogIndex, prevLogTerm := insertLogTable(term, command, 1)
	logFile("append", "id: "+strconv.Itoa(prevLogIndex))
	s := getState()
	message := myNodeID + " " + AppendEntryRPC + " " + strconv.Itoa(s.currentTerm) + " " + strconv.Itoa(prevLogIndex) +
		" " + strconv.Itoa(prevLogTerm) + " " + command + " " + strconv.Itoa(s.commitIndex) + "\n"
	logFile("append", message)
	for node := range otherNodes {
		chanMap[otherNodes[node].nodeID] <- message
	}
	logFile("append", "Append Entry Init Ends\n")
}

func handleAppendEntryRPCFromLeader(message string) {
	logFile("append", "Append Entry Starts\n")

	logFile("append", "Append Entry Ends\n")
}
