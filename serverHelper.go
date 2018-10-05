package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func initiateElection(myNodeID string) {

	fmt.Println("initiate Election Begins:", time.Now())
	term, _ := getState()
	insertTableState(term+1, myNodeID)
	message := myNodeID + " " + RequestVoteRPC + " " + strconv.Itoa(term+1) + "\n"
	fmt.Println("conn map len:", len(connMap))
	for rID := range otherNodes {
		go sendMessage(message, otherNodes[rID].nodeID)
	}
	fmt.Println("initiate Election Ends")
}

func processRequest(message string, conn net.Conn) {
	fmt.Println("process Request Starts")
	data := strings.Fields(message)
	fmt.Println(data)
	if data[1] == RequestVoteRPC {
		handleRequestVoteRPC(data, conn, data[0])
	}
	fmt.Println("process Request Ends")
}

func handleRequestVoteRPC(data []string, conn net.Conn, remoteNodeID string) {
	fmt.Println("handleRequestVoteRPC starts")
	candidateTerm, _ := strconv.Atoi(data[2])
	myTerm, votedFor := getState()
	var message string
	fmt.Println("myTerm:", myTerm, " votedFor: ", votedFor)
	if votedFor == "" && candidateTerm > myTerm {
		message = myNodeID + " " + RequestVoteRPCReply + " " + "YES\n"
		res := insertTableState(candidateTerm, remoteNodeID)
		fmt.Println("res insert status: ", res)
	} else {
		message = myNodeID + " " + RequestVoteRPCReply + " " + "NO" + " " + strconv.Itoa(myTerm) + "\n"
	}
	fmt.Println("handleRequestVoteRPC:", message)
	go sendMessage(message, remoteNodeID)
	//_, err := conn.Write([]byte(message))
	//checkError(err, "handleRequestVoteRPC")
	resetTimer()
	fmt.Println("handleRequestVoteRPC ends")
}
