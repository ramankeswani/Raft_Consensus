package main

import (
	"fmt"
	"strconv"
)

/*
Creates Append Entry RPC Request and sends the request to all Clients
*/
func appendEntryInit(command string) {
	fmt.Println("Append Entry Starts Map len:", len(connMap))
	prevLogIndex, prevLogTerm := insertLogTable(term, command, 1)
	fmt.Println("id:", prevLogIndex)
	s := getState()
	message := myNodeID + " " + AppendEntryRPC + " " + strconv.Itoa(s.currentTerm) + " " + strconv.Itoa(prevLogIndex) +
		" " + strconv.Itoa(prevLogTerm) + " " + command + " " + strconv.Itoa(s.commitIndex)
	for node := range otherNodes {
		chanMap[otherNodes[node].nodeID] <- message
	}
	fmt.Println("Append Entry Ends")
}

func handleAppendEntryRPCFromLeader(message string) {
	fmt.Println("???handle Append Entry RPC from leader starts???")

	fmt.Println("???handle Append Entry RPC from leader starts???")
}
