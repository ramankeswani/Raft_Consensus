package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

/*
Invoked when a node restarts after failure. Communicates with other nodes to establish socket channels for communication
*/
func initRecovery(port string) {
	logFile("recover", "initRecoveryStarts\n")
	message := nodeID + " " + NodeRecoverMessage + " " + port + "\n"
	for node := range otherNodes {
		fmt.Println("init Recovery: chanMap[otherNodes[node].nodeID] nil?: " + strconv.FormatBool(nil != chanMap[otherNodes[node].nodeID]))
		if nil != chanMap[otherNodes[node].nodeID] {
			chanMap[otherNodes[node].nodeID] <- message
		}
	}
	logFile("recover", "initRecoveryEnds\n")
}

/*
Active Recevies Recovery Message from failed node to re-esatablish client - server socket channel
*/
func handleRecoveryMessage(message string) {
	logFile("recover", "handleRecoveryMessage Starts\n")
	dataSlice := strings.Split(strings.TrimRight(message, "\n"), " ")
	port, _ := strconv.Atoi(dataSlice[2])
	ip := getIP(dataSlice[0])
	address := ip + ":" + strconv.Itoa(port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	checkError(err, "client")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	c := connection{
		nodeID: dataSlice[0],
		conn:   conn,
	}
	connMap[dataSlice[0]] = c
	//totalNodes++
	logFile("recover", "handleRecoveryMessage Ends\n")
}

/*
Invoked when a follower rejects AppendEntryRequest.
Decrements NextIndex and send AppendEntry until follower accepts.
*/
func synchronizeLogs(nodeID string, logIndex int) {
	logFile("recover", "synchronize logs starts nodeid: "+nodeID+" logIndex: "+strconv.Itoa(logIndex)+"\n")
	nextIndex[nodeID] = logIndex - 1
	go sendAppendRequest(nodeID, nextIndex[nodeID])
	logFile("recover", "synchronize logs ends\n")
}

/*
Prepares message for Append Entry
*/
func prepareAppendEntryRequest(nodeID string, logIndex int) (message string) {
	logFile("recover", "prepareAppendEntryRequest starts node id: "+nodeID+" logID: "+strconv.Itoa(logIndex)+"\n")
	s := getState()
	prevLogIndex := "-1"
	prevLogTerm := "-1"
	commitFlag := "0"
	if logIndex > 1 {
		prevLog := getLogTable(logIndex - 1)
		prevLogIndex = strconv.Itoa(prevLog.logIndex)
		prevLogTerm = strconv.Itoa(prevLog.term)
	}
	log := getLogTable(logIndex)
	commitFlag = strconv.Itoa(log.isCommited)
	message = myNodeID + " " + SyncRequest + " " + strconv.Itoa(s.currentTerm) + " " + prevLogIndex + " " +
		prevLogTerm + " " + log.command + " " + strconv.Itoa(logIndex) + " " + commitFlag + "\n"
	logFile("recover", "prepareAppendEntryRequest ends \n")
	return message
}

/*
Sends Append Entry to Node being recovered
*/
func sendAppendRequest(nodeID string, logIndex int) {
	logFile("recover", "sendAppendRequest starts \n")
	// LeaderNodeID | SyncRequest | Current Term | PrevLogIndex | PrevLogTerm | Log Command | New Log Index | CommitFlag
	message := prepareAppendEntryRequest(nodeID, logIndex)
	chanMap[nodeID] <- message
	response := <-chanSyncResp
	logFile("recover", "sendAppendRequest respones: "+response+"\n")
	dataSlice := strings.Split(strings.TrimSuffix(response, "\n"), " ")
	if strings.Compare(dataSlice[2], ACCEPT) == 0 {
		//go handleAppendEntryRPCReply(response)
		if strings.Compare(dataSlice[3], "1") == 0 {
			nextIndex[nodeID] = logIndex + 1
		}
		go overwriteLogs(nodeID, logIndex+1)
	}
	/* else {
		go sendAppendRequest(nodeID, logIndex-1)
	} */
	logFile("recover", "sendAppendRequest Ends \n")
}

/*
Overwrite logs to the follower
*/
func overwriteLogs(nodeID string, logIndex int) {
	logFile("recover", "overwriteLogs Starts \n")
	n := getLatestLog().logIndex
	if logIndex <= n {
		go sendAppendRequest(nodeID, logIndex)
	}
	logFile("recover", "overwriteLogs Ends \n")
}

/*
Handles SyncRequest and sends back appropriate response to leader
*/
func handleSyncRequest(message string) {
	logFile("recover", "handleSyncRequest Starts \n")
	go handleAppendEntryRPCFromLeader(message, true)
	response := <-chanAppendResp
	logFile("recover", "handleSyncRequest response: "+response+"\n")
	dataSlice := strings.Split(strings.TrimRight(response, "\n"), " ")
	message = dataSlice[0] + " " + SyncRequestReply + " " + dataSlice[2] + " " + dataSlice[3] + "\n"
	logFile("recover", "handleSyncRequest message: "+message)
	s := getState()
	go temp(message, dataSlice[0], s.leader)
	logFile("recover", "handleSyncRequest Ends \n")
}

func temp(message string, nodeID string, leader string) {
	logFile("recover", "temp leader: "+leader+" nodeid: "+nodeID+" null == chanMap[leadar] "+strconv.FormatBool(chanMap[leader] == nil))
	for k := range chanMap {
		logFile("recover", "map key: "+k+"\n")
	}
	chanMap[leader] <- message
	logFile("recover", "temp Ends \n")
}

/*
Initiates log resync for recovering node. Invoked upon receving marker - SyncOnLoad
*/
func initLogSync(message string) {
	logFile("recover", "initLogSync Starts \n")
	dataSlice := strings.Split(strings.TrimRight(message, "\n"), " ")
	log := getLatestLog()
	fmt.Println("init Log Sync log index: " + strconv.Itoa(log.logIndex))
	if log.logIndex >= 1 {
		go synchronizeLogs(dataSlice[0], log.logIndex+1)
	}
	logFile("recover", "initLogSync Ends \n")
}
