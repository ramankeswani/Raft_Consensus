package main

import (
	"strconv"
	"strings"
	"sync"
)

var mutexUpdateVote = &sync.Mutex{}
var commitIndex []int

/*
Creates Append Entry RPC Request and sends the request to all Clients
*/
func appendEntryInit(command string) {
	logFile("append", "Append Entry Init Starts Map len:"+strconv.Itoa(len(connMap))+" command: "+command+"\n")
	prevLogIndex, prevLogTerm := insertLogTable(term, command, 1)
	logFile("append", "id: "+strconv.Itoa(prevLogIndex)+"\n")
	s := getState()
	// LeaderNodeID | AppendEntryRPC | Current Term | PrevLogIndex | PrevLogTerm | Log Command | New Log Index | CommitIndex
	tempCurrentIndex := prevLogIndex + 1
	if prevLogIndex == -1 {
		tempCurrentIndex = 1
	}
	message := myNodeID + " " + AppendEntryRPC + " " + strconv.Itoa(s.currentTerm) + " " + strconv.Itoa(prevLogIndex) +
		" " + strconv.Itoa(prevLogTerm) + " " + command + " " + strconv.Itoa(tempCurrentIndex) + " " + strconv.Itoa(s.commitIndex) + "\n"
	logFile("append", message)
	for node := range otherNodes {
		chanMap[otherNodes[node].nodeID] <- message
	}
	logFile("append", "Append Entry Init Ends\n")
}

/*
Invoked when a follower receives Append Entry RPC from Leader
*/
func handleAppendEntryRPCFromLeader(message string) {
	logFile("append", "handleAppendEntryRPCFromLeader Starts\n")
	s := getState()
	dataSlice := strings.Split(message, " ")
	rpcTerm, _ := strconv.Atoi(dataSlice[2])
	logFile("append", "state curr term: "+strconv.Itoa(s.currentTerm)+" rpcterm: "+dataSlice[2]+
		" dataslice[0]: "+dataSlice[0]+" New Log Index: "+dataSlice[6]+"\n")
	if s.currentTerm == rpcTerm {
		pLID, pLT := insertLogTable(rpcTerm, dataSlice[5], 0)
		logFile("append", "handleAppendEntryRPCFromLeader message: "+myNodeID+" "+AppendEntryRPCReply+" YES "+
			"prevIndex: "+strconv.Itoa(pLID)+" pLogTerm: "+strconv.Itoa(pLT)+" "+dataSlice[6]+"\n")
		chanMap[dataSlice[0]] <- myNodeID + " " + AppendEntryRPCReply + " YES " + dataSlice[6] + "\n"
	}
	logFile("append", "handleAppendEntryRPCFromLeader Ends\n")
}

func handleAppendEntryRPCReply(message string) {
	mutexUpdateVote.Lock()
	logFile("append", "handleAppendEntryRPCReply Starts\n")
	dataSlice := strings.Split(strings.TrimSuffix(message, "\n"), " ")
	logFile("append", "handleAppendEntryRPCReply message: "+message)
	if strings.Compare(dataSlice[2], YES) == 0 {
		logIndex, _ := strconv.Atoi(dataSlice[3])
		logFile("append", "handleAppendEntryRPCReply dataslice[3]: "+dataSlice[3]+" logindex: "+strconv.Itoa(logIndex)+"\n")
		votes := incrementVoteCount(logIndex)
		logFile("append", "handleAppendEntryRPCReply votes: "+strconv.Itoa(votes)+"\n")
		if float32(float32(votes)/float32(totalNodes)) > 0.5 {
			logFile("append", "majority\n")
		}
	}
	logFile("append", "handleAppendEntryRPCReply Ends\n")
	mutexUpdateVote.Unlock()
}

/*
Leader sends out commit request when a log is replicated on majority of followers
*/
func sendCommitRequest(logIndex int) {
	logFile("commit", "sendCommitRequest Starts logIndex: "+strconv.Itoa(logIndex))
	prevLogEntry := getLogTable(logIndex - 1)
	s := getState()
	message := myNodeID + " " + strconv.Itoa(s.currentTerm) + " CommitEntryRequest " +
		strconv.Itoa(prevLogEntry.logIndex) + " " + strconv.Itoa(prevLogEntry.term) + " " + strconv.Itoa(logIndex) + "\n"
	logFile("commit", "sendCommitRequest Message TEMP:"+message)
	for node := range otherNodes {
		nextID := nextIndex[otherNodes[node].nodeID]
		message = myNodeID + " " + strconv.Itoa(s.currentTerm) + " CommitEntryRequest " +
			strconv.Itoa(prevLogEntry.logIndex) + " " + strconv.Itoa(prevLogEntry.term) + " " +
			strconv.Itoa(nextID) + "\n"
		chanMap[otherNodes[node].nodeID] <- message
		//TO-DO Update nextID
	}
}

/*
Initialized next Commit Index Map - Stores next entry to be commited by Followe
Invoked when a candidate becomes leader
Initial value is the highest log index in leader's own log
*/
func initNextIndexMap() (nextIndex map[string]int) {
	logFile("commit", "initNextIndexMap starts\n")
	nextIndex = make(map[string]int)
	index := getLatestLog()
	if index == -1 {
		index = 1
	}
	for node := range otherNodes {
		nextIndex[otherNodes[node].nodeID] = index
	}
	logFile("commit", "Iterating Map\n")
	for k, v := range nextIndex {
		logFile("commit", "k: "+k+" v: "+strconv.Itoa(v)+"\n")
	}
	logFile("commit", "initNextIndexMap ends\n")
	return nextIndex
}
