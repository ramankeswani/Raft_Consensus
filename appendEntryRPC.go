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

// TO-DO - Move to persistent
var isCommited map[int]bool

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
			logFile("commit", "majority isCommited[logIndex]: "+strconv.FormatBool(isCommited[logIndex])+"\n")
			if !isCommited[logIndex] {
				sendCommitRequest(logIndex)
			}
			isCommited[logIndex] = true
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
	message := myNodeID + " " + strconv.Itoa(s.currentTerm) + " " + CommitEntryRequest + " " +
		strconv.Itoa(prevLogEntry.logIndex) + " " + strconv.Itoa(prevLogEntry.term) + " " + strconv.Itoa(logIndex) + "\n"
	logFile("commit", "sendCommitRequest Message TEMP:"+message)
	for node := range otherNodes {
		nextID := nextIndex[otherNodes[node].nodeID]
		// myNodeID |CommitEntryRequest |  Current Term | PrevLogIndex | PrevLogTerm | Next Commit Index
		message = myNodeID + " " + CommitEntryRequest + " " + strconv.Itoa(s.currentTerm) + " " +
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
	index := getLatestLog().logIndex
	if index < 0 {
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

/*
Invoked by client's server to process Commit Entry Request from Leader
Sends a confirmation back to leader if prev Log entry matches(and commits the log) or rejects otherwise
Input Message Format:  LeaderNodeID | CommitEntryRequest | Current Term | PrevLogIndex | PrevLogTerm | Next Commit Index
*/
func handleCommitEntryRequest(message string) {
	logFile("commit", "handleCommitEntryRequest Starts\n")
	s := getState()
	message = strings.TrimSuffix(message, "\n")
	dataSlice := strings.Split(message, " ")
	leaderTerm, _ := strconv.Atoi(dataSlice[2])
	logIndex, _ := strconv.Atoi(dataSlice[5])
	logFile("commit", "handleCommitEntryRequest leader term : "+dataSlice[2]+" logIndex: "+dataSlice[5]+
		" leaderid: "+dataSlice[0]+" myterm: "+strconv.Itoa(s.currentTerm)+"\n")
	if s.currentTerm > leaderTerm {
		message = myNodeID + " " + strconv.Itoa(s.currentTerm) + " " + s.leader + " " + REJECT + "\n"
		chanMap[dataSlice[0]] <- message
		return
	}
	prevLogIndex, _ := strconv.Atoi(dataSlice[3])
	prevLogTerm, _ := strconv.Atoi(dataSlice[4])
	if prevLogIndex != -1 {
		if !appendRPCCheck(prevLogIndex, prevLogTerm) {
			message = myNodeID + " " + dataSlice[2] + " " + s.leader + " " + REJECT + "\n"
			chanMap[dataSlice[0]] <- message
			return
		}
	}
	commitStatus := commitLog(logIndex)
	if commitStatus {
		message = myNodeID + " " + dataSlice[2] + " " + s.leader + " " + ACCEPT + "\n"
	} else {
		message = myNodeID + " " + dataSlice[2] + " " + s.leader + " " + REJECT + "\n"
	}
	logFile("commit", "handleCommitEntryRequest message:"+message)
	logFile("commit", "handleCommitEntryRequest chan:"+strconv.FormatBool(chanMap[dataSlice[0]] != nil)+"\n")
	chanMap[dataSlice[0]] <- message
	logFile("commit", "handleCommitEntryRequest Ends\n")
}

/*
Compares Previous Log Index and Term to ensure safe replication of logs
*/
func appendRPCCheck(prevLogIndex int, prevLogTerm int) (matches bool) {
	logFile("commit", "appendRPCCheck Starts\n")
	l := getLatestLog()
	l = getLogTable(l.logIndex - 1)
	matches = true
	if l.logIndex != prevLogIndex || l.term != prevLogTerm {
		matches = false
	}
	logFile("commit", "appendRPCCheck Ends matches: "+strconv.FormatBool(matches)+"\n")
	return matches
}
