package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var candidate bool
var leader bool
var votes int
var nextIndex map[string]int

/*
Sends RequestVoteRPC to all other Nodes
Arguments: Self Node ID
*/
func initiateElection(myNodeID string) {

	fmt.Println("initiate Election Begins:", time.Now())
	candidate = true
	votes = 0
	s := getState()
	term = s.currentTerm + 1
	insertTableState(term, myNodeID, "", 0)
	message := myNodeID + " " + RequestVoteRPC + " " + strconv.Itoa(s.currentTerm+1) + "\n"
	fmt.Println("election messsage:", message)
	for remoteID := range otherNodes {
		chanMap[otherNodes[remoteID].nodeID] <- message
	}
	fmt.Println("initiate Election Ends")
}

/*
Scans all incoming messages to find type of message
Arguments: Message from remote
*/
func processRequest(message string) {
	//fmt.Println("process Request Starts")
	data := strings.Fields(message)
	fmt.Println(data)
	if data[1] == RequestVoteRPC {
		handleRequestVoteRPC(data, data[0])
	} else if data[1] == RequestVoteRPCReply && candidate {
		countVotes(data)
	} else if data[1] == AppendEntryFromClient {
		appendEntryInit(data[2])
	} else if data[1] == AppendEntryRPC {
		handleAppendEntryRPCFromLeader(message)
	} else if data[1] == AppendEntryRPCReply {
		handleAppendEntryRPCReply(message)
	} else if data[1] == CommitEntryRequest {
		handleCommitEntryRequest(message)
	}
	//fmt.Println("process Request Ends")
}

/*
Keep tracks of Votes Received
*/
func countVotes(data []string) {
	fmt.Println("Count Votes Start: ", data[2])
	if strings.Compare(YES, data[2]) == 0 {
		votes++
		fmt.Printf("Votes: %d totalNodes:  %d \n", votes, totalNodes)
		fmt.Printf("votes/totalNodes: %f \n", float32(float32(votes)/float32(totalNodes)))
		if float32(float32(votes)/float32(totalNodes)) > 0.5 {
			killTimer()
			s := getState()
			insertTableState(s.currentTerm, s.votedFor, myNodeID, 0)
			leader = true
			term = s.currentTerm
			go heartbeat(otherNodes, myNodeID, connMap, getState())
			candidate = false
			votes = 0
			// Stores Next Commit Index for each Follower
			nextIndex = initNextIndexMap()
		}
	}
	fmt.Println("Count Votes END")
}

/*
Determines if Vote should be granted to Candidate
Arguments: Message from remote, remote NodeID
*/
func handleRequestVoteRPC(data []string, remoteNodeID string) {
	fmt.Println("handleRequestVoteRPC starts")
	go resetTimer()
	candidateTerm, _ := strconv.Atoi(data[2])
	s := getState()
	var message string
	fmt.Println("From DB: myTerm:", s.currentTerm, " votedFor: ", s.votedFor)
	if candidateTerm > s.currentTerm {
		message = myNodeID + " " + RequestVoteRPCReply + " " + "YES\n"
		res := insertTableState(candidateTerm, remoteNodeID, "", 0)
		leader = false
		candidate = false
		term = candidateTerm
		fmt.Println("res insert status: ", res)
	} else {
		message = myNodeID + " " + RequestVoteRPCReply + " " + "NO" + " " + strconv.Itoa(s.currentTerm) + "\n"
	}
	fmt.Println("handleRequestVoteRPC:", message)
	chanMap[remoteNodeID] <- message
	//go sendMessage(message, remoteNodeID)
	fmt.Println("handleRequestVoteRPC ends")
}
