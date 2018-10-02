package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func initiateElection(nodeID string) {

	fmt.Println("initiate Election Begins:", time.Now())
	term, _ := getState()
	insertTableState(term+1, nodeID)
	message := RequestVoteRPC + " " + strconv.Itoa(term+1) + " " + nodeID
	for conn := range connections {
		connections[conn].Write([]byte(message))
	}
	fmt.Println("initiate Election Ends")
}

func processRequest(message string, conn net.Conn) {
	fmt.Println("process Request Starts")
	data := strings.Fields(message)
	fmt.Println(data)
	if data[0] == RequestVoteRPC {
		handleRequestVoteRPC(data, conn)
	}
	fmt.Println("process Request Ends")
}

func handleRequestVoteRPC(data []string, conn net.Conn) {
	candidateTerm, _ := strconv.Atoi(data[1])
	myTerm, votedFor := getState()
	var message string
	fmt.Println("myTerm:", myTerm, " votedFor: ", votedFor)
	if votedFor == "" && candidateTerm > myTerm {
		message = RequestVoteRPCReply + " " + "YES"
		res := insertTableState(myTerm, data[2])
		fmt.Println("res insert status: ", res)
	} else {
		message = RequestVoteRPCReply + " " + "NO" + " " + strconv.Itoa(myTerm)
	}
	fmt.Println("handleRequestVoteRPC:", message)
	_, err := conn.Write([]byte(message))
	checkError(err, "handleRequestVoteRPC")
}
