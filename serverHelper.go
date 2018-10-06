package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

/*
Sends RequestVoteRPC to all other Nodes
Arguments: Self Node ID
*/
func initiateElection(myNodeID string) {

	fmt.Println("initiate Election Begins:", time.Now())
	term, _ := getState()
	insertTableState(term+1, myNodeID)
	message := myNodeID + " " + RequestVoteRPC + " " + strconv.Itoa(term+1) + "\n"
	for rID := range otherNodes {
		go sendMessage(message, otherNodes[rID].nodeID)
	}
	fmt.Println("initiate Election Ends")
}

/*
Scans all incoming messages to find type of message
Arguments: Message from remote
*/
func processRequest(message string) {
	fmt.Println("process Request Starts")
	data := strings.Fields(message)
	fmt.Println(data)
	if data[1] == RequestVoteRPC {
		handleRequestVoteRPC(data, data[0])
	}
	fmt.Println("process Request Ends")
}

/*
Determines if Vote should be granted to Candidate
Arguments: Message from remote, remote NodeID
*/
func handleRequestVoteRPC(data []string, remoteNodeID string) {
	fmt.Println("handleRequestVoteRPC starts")
	candidateTerm, _ := strconv.Atoi(data[2])
	myTerm, votedFor := getState()
	var message string
	fmt.Println("From DB: myTerm:", myTerm, " votedFor: ", votedFor)
	if votedFor == "" && candidateTerm > myTerm {
		message = myNodeID + " " + RequestVoteRPCReply + " " + "YES\n"
		res := insertTableState(candidateTerm, remoteNodeID)
		fmt.Println("res insert status: ", res)
	} else {
		message = myNodeID + " " + RequestVoteRPCReply + " " + "NO" + " " + strconv.Itoa(myTerm) + "\n"
	}
	fmt.Println("handleRequestVoteRPC:", message)
	go sendMessage(message, remoteNodeID)
	resetTimer()
	fmt.Println("handleRequestVoteRPC ends")
}
