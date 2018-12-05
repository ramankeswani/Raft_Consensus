package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mutexUpdateVote = &sync.Mutex{}
var mutexAppendRPC = &sync.Mutex{}
var commitIndex []int
var startTime time.Time
var dif time.Duration

/*
Creates Append Entry RPC Request and sends the request to all Clients
*/
func appendEntryInit(command string) {
	mutexAppendRPC.Lock()
	startTime = time.Now()
	logFile("append", "Append Entry Init Starts Map len:"+strconv.Itoa(len(connMap))+" command: "+command+"\n")
	prevLogIndex, prevLogTerm := insertLogTable(term, command, 1)
	logFile("append", "id: "+strconv.Itoa(prevLogIndex)+"\n")
	s := getState()
	// LeaderNodeID | AppendEntryRPC | Current Term | PrevLogIndex | PrevLogTerm | Log Command | New Log Index | NextIndex | CommitFLag
	tempCurrentIndex := prevLogIndex + 1
	if prevLogIndex == -1 {
		tempCurrentIndex = 1
	}
	baseMessage := myNodeID + " " + AppendEntryRPC + " " + strconv.Itoa(s.currentTerm) + " " + strconv.Itoa(prevLogIndex) +
		" " + strconv.Itoa(prevLogTerm) + " " + command + " " + strconv.Itoa(tempCurrentIndex)
	logFile("append", baseMessage)
	for node := range otherNodes {
		message := baseMessage + " " + strconv.Itoa(nextIndex[otherNodes[node].nodeID]) + " 0\n"
		logFile("commit", "AppendEntryInit message: "+message)
		chanMap[otherNodes[node].nodeID] <- message
	}
	logFile("append", "Append Entry Init Ends\n")
	mutexAppendRPC.Unlock()
}

/*
Invoked when a follower receives Append Entry RPC from Leader
Incoming Message Format:
LeaderNodeID | AppendEntryRPC | Current Term | PrevLogIndex | PrevLogTerm | Log Command | New Log Index | NextIndex | CommitFlag
*/
func handleAppendEntryRPCFromLeader(message string, sendRespFlag bool) {
	logFile("recover", "handleAppendEntryRPCFromLeader Starts flag: "+strconv.FormatBool(sendRespFlag)+"\n")
	s := getState()
	dataSlice := strings.Split(strings.TrimRight(message, "\n"), " ")
	rpcTerm, _ := strconv.Atoi(dataSlice[2])
	fmt.Println(dataSlice)
	logFile("recover", "state curr term: "+strconv.Itoa(s.currentTerm)+" rpcterm: "+dataSlice[2]+
		" dataslice[0]: "+dataSlice[0]+" New Log Index: "+dataSlice[6]+"\n")
	if s.currentTerm == rpcTerm {
		prevLogIndex, _ := strconv.Atoi(dataSlice[3])
		prevLogTerm, _ := strconv.Atoi(dataSlice[4])
		if prevLogIndex != -1 {
			if !appendRPCCheck(prevLogIndex, prevLogTerm) {
				message = myNodeID + " " + AppendEntryRPCReply + " " + REJECT + " " + dataSlice[2] + " " + s.leader +
					" " + dataSlice[6] + "\n"
				chanMap[dataSlice[0]] <- message
				if sendRespFlag {
					logFile("recover", "handleAppendEntryRPCFromLeader sendRespFlag true\n")
					chanAppendResp <- message
					logFile("recover", "handleAppendEntryRPCFromLeader Ends reject\n")
				}
				return
			}
		}
		pLID, pLT := insertLogTable(rpcTerm, dataSlice[5], 0)
		logFile("recover", "handleAppendEntryRPCFromLeader message: "+myNodeID+" "+AppendEntryRPCReply+" YES "+
			"prevIndex: "+strconv.Itoa(pLID)+" pLogTerm: "+strconv.Itoa(pLT)+" newlogindex: "+dataSlice[6]+"\n")
		message = myNodeID + " " + AppendEntryRPCReply + " " + ACCEPT + " " + dataSlice[6] + "\n"
		chanMap[dataSlice[0]] <- message
		if sendRespFlag {
			logFile("recover", "chanAppendResp <- "+message)
			logFile("recover", "null == chanAppendResp "+strconv.FormatBool(nil == chanAppendResp)+"\n")
			chanAppendResp <- message
			if strings.Compare(dataSlice[7], "1") == 0 {
				commitLog(pLID + 1)
			}
			logFile("recover", "message sent into channel\n")
		}
	}
	logFile("recover", "handleAppendEntryRPCFromLeader Ends\n")
}

// TO-DO - Move to persistent
//var isCommited map[int]bool

/*
Incoming Message Format:
NodeID | AppendEntryRPCReply | REJECT | Term | Leader | LogIndex
NodeID | AppendEntryRPCReply | ACCEPT | LogIndex
*/
func handleAppendEntryRPCReply(message string) {
	mutexUpdateVote.Lock()
	logFile("recover", "handleAppendEntryRPCReply Starts\n")
	dataSliceT := strings.Split(strings.TrimSuffix(message, "\n"), " ")
	logFile("recover", "handleAppendEntryRPCReply message: "+message)
	if strings.Compare(dataSliceT[2], ACCEPT) == 0 {
		lI, _ := strconv.Atoi(dataSliceT[3])
		logFile("recover", "handleAppendEntryRPCReply term]: "+dataSliceT[3]+" logindex: "+strconv.Itoa(lI)+"\n")
		votes := incrementVoteCount(lI)
		logFile("recover", "handleAppendEntryRPCReply votes: "+strconv.Itoa(votes)+"\n")
		if float32(float32(votes)/float32(totalNodes)) > 0.5 {
			l := getLogTable(lI)
			logFile("recover", "majority isCommited[logIndex]: "+strconv.Itoa(l.isCommited)+"\n")
			if l.isCommited == 0 {
				sendCommitRequest(lI)
				commitLog(lI)
			}
		}
	} else {
		logIndexT, _ := strconv.Atoi(dataSliceT[5])
		go synchronizeLogs(dataSliceT[0], logIndexT)
	}
	logFile("recover", "handleAppendEntryRPCReply Ends\n")
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
		index = 0
	}
	for node := range otherNodes {
		nextIndex[otherNodes[node].nodeID] = index + 1
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
Input Message Format:
LeaderNodeID | CommitEntryRequest | Current Term | PrevLogIndex | PrevLogTerm | Next Commit Index
*/
func handleCommitEntryRequest(message string) {
	logFile("commit", "handleCommitEntryRequest Starts\n")
	s := getState()
	message = strings.TrimSuffix(message, "\n")
	dataSlice := strings.Split(message, " ")
	logIndex, _ := strconv.Atoi(dataSlice[5])
	logFile("commit", "handleCommitEntryRequest leader term : "+dataSlice[2]+" logIndex: "+dataSlice[5]+
		" leaderid: "+dataSlice[0]+" myterm: "+strconv.Itoa(s.currentTerm)+"\n")
	commitStatus := commitLog(logIndex)
	if commitStatus {
		message = myNodeID + " " + CommitEntryReply + " " + dataSlice[2] + " " + s.leader + " " + ACCEPT + " " + strconv.Itoa(logIndex) + "\n"
	} else {
		message = myNodeID + " " + CommitEntryReply + " " + dataSlice[2] + " " + s.leader + " " + REJECT + " " + strconv.Itoa(logIndex) + "\n"
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
	matches = true
	if l.logIndex != prevLogIndex || l.term != prevLogTerm {
		matches = false
	}
	logFile("commit", "appendRPCCheck Ends matches: "+strconv.FormatBool(matches)+"\n")
	return matches
}

/*
Invoked when Followers reply to Commit Entry Request
Updates the nextIndex for follower
Input Message Format:
FollowerID | CommitEntryReply | Current Term | Leader ID | Response | LogIndex
*/
func handleCommitEntryReply(message string) {
	dataSlice := strings.Split(strings.TrimRight(message, "\n"), " ")
	logFile("commit", "handleCommitEntryReply Starts data slice 5: "+dataSlice[5]+"\n")
	tempLogIndex, _ := strconv.Atoi(dataSlice[5])
	if strings.Compare(dataSlice[4], ACCEPT) == 0 {
		logFile("commit", "logIndex: "+strconv.Itoa(tempLogIndex)+" "+"map: key: "+dataSlice[0]+" val: "+strconv.Itoa(nextIndex[dataSlice[0]])+"\n")
		nextIndex[dataSlice[0]] = tempLogIndex + 1
		logFile("commit", "Next Index "+dataSlice[0]+" "+strconv.Itoa(tempLogIndex+1)+"\n")
		//if dif.Nanoseconds() != 0 {
		dif = time.Now().Sub(startTime)
		fmt.Println("----------replication timeout-----------: " + strconv.Itoa(int(dif)))
		//}
	}
	logFile("commit", "handleCommitEntryReply Ends\n")
}
