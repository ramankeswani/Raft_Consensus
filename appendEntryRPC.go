package main

import (
	"strconv"
	"strings"
)

/*
Creates Append Entry RPC Request and sends the request to all Clients
*/
func appendEntryInit(command string) {
	logFile("append", "Append Entry Starts Map len:"+strconv.Itoa(len(connMap))+" command: "+command+"\n")
	prevLogIndex, prevLogTerm := insertLogTable(term, command, 1)
	logFile("append", "id: "+strconv.Itoa(prevLogIndex)+"\n")
	s := getState()
	// LeaderNodeID | AppendEntryRPC | Current Term | PrevLogIndex | PrevLogTerm | Log Command | New Log Index | CommitIndex
	message := myNodeID + " " + AppendEntryRPC + " " + strconv.Itoa(s.currentTerm) + " " + strconv.Itoa(prevLogIndex) +
		" " + strconv.Itoa(prevLogTerm) + " " + command + " " + strconv.Itoa(prevLogIndex+1) + " " + strconv.Itoa(s.commitIndex) + "\n"
	logFile("append", message)
	for node := range otherNodes {
		chanMap[otherNodes[node].nodeID] <- message
	}
	logFile("append", "Append Entry Init Ends\n")
}

func handleAppendEntryRPCFromLeader(message string) {
	logFile("append", "Append Entry Starts\n")
	s := getState()
	dataSlice := strings.Split(message, " ")
	rpcTerm, _ := strconv.Atoi(dataSlice[2])
	logFile("append", "state curr term: "+strconv.Itoa(s.currentTerm)+" rpcterm: "+dataSlice[2]+
		" dataslice[0]: "+dataSlice[0]+" commitIndex: "+dataSlice[6]+"\n")
	if s.currentTerm == rpcTerm {
		insertLogTable(rpcTerm, dataSlice[5], 0)
		chanMap[dataSlice[0]] <- myNodeID + " " + AppendEntryRPCReply + " YES " + dataSlice[6] + "\n"
	}
	logFile("append", "Append Entry Ends\n")
}

func handleAppendEntryRPCReply(message string) {
	logFile("append", "handleAppendEntryRPCReply Starts\n")
	dataSlice := strings.Split(message, " ")
	if strings.Compare(dataSlice[2], YES) == 0 {
		logIndex, _ := strconv.Atoi(dataSlice[3])
		votes := incrementVoteCount(logIndex)
		logFile("append", "handleAppendEntryRPCReply votes: "+strconv.Itoa(votes)+"\n")
		if float32(float32(votes)/float32(totalNodes)) > 0.5 {
			logFile("append", "majority\n")
		}
	}
	logFile("append", "handleAppendEntryRPCReply Ends\n")
}
