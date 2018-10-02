package main

import (
	"net"
	"strconv"
)

func requestVote(candidateTerm int, candidateID string, conn *net.TCPConn) {

	term := strconv.Itoa(candidateTerm + 1)
	_, err := conn.Write([]byte(RequestVoteRPC + " " + term + " " + candidateID))
	checkError(err, "RequestVoteRPC")

}
