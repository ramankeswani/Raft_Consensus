package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var candidate bool
var votes int
var connMapServerHelper map[string]connection

/*
Sends RequestVoteRPC to all other Nodes
Arguments: Self Node ID
*/
func initiateElection(myNodeID string) {

	fmt.Println("initiate Election Begins:", time.Now())
	candidate = true
	votes = 0
	s := getState()
	insertTableState(s.currentTerm+1, myNodeID, "")
	message := myNodeID + " " + RequestVoteRPC + " " + strconv.Itoa(s.currentTerm+1) + "\n"
	fmt.Println("election messsage:", message)
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
	} else if data[1] == RequestVoteRPCReply {
		countVotes(data)
	}
	fmt.Println("process Request Ends")
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
			s := getState()
			insertTableState(s.currentTerm, s.votedFor, myNodeID)
			heartbeat(otherNodes, myNodeID, connMapServerHelper, s)
			candidate = false
			votes = 0
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
	candidateTerm, _ := strconv.Atoi(data[2])
	s := getState()
	var message string
	fmt.Println("From DB: myTerm:", s.currentTerm, " votedFor: ", s.votedFor)
	if candidateTerm > s.currentTerm {
		message = myNodeID + " " + RequestVoteRPCReply + " " + "YES\n"
		res := insertTableState(candidateTerm, remoteNodeID, "")
		fmt.Println("res insert status: ", res)
	} else {
		message = myNodeID + " " + RequestVoteRPCReply + " " + "NO" + " " + strconv.Itoa(s.currentTerm) + "\n"
	}
	fmt.Println("handleRequestVoteRPC:", message)
	go sendMessage(message, remoteNodeID)
	go resetTimer()
	fmt.Println("handleRequestVoteRPC ends")
}

/*
Initialize connection map for server helper
*/
func initConnectionMapServerHelper(chanConnMap chan map[string]connection) {
	connMapServerHelper = <-chanConnMap
}
