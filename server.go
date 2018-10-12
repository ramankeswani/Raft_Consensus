package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

var timer *time.Timer
var source rand.Source
var r *rand.Rand
var myNodeID string
var term int

/*
Server Modules Invoked from Main on load
Arguements: Server Port, Self Node ID
Returns/Exits on Node Failure Only
*/
func server(myPort int, nodeID string) {

	source = rand.NewSource(time.Now().UnixNano())
	r = rand.New(source)
	myNodeID = nodeID
	fmt.Println("mynode id", myNodeID)
	service := ":" + strconv.Itoa(myPort)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err, "server")
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err, "server")

	// Routine to track HeartBeat Timeout
	go heartbeatChecker()

	for {
		fmt.Println("Server Ready to Accept")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error at server while accepting")
			continue
		}
		fmt.Println("Server Accepted From: ", conn.RemoteAddr())

		go handleRequests(conn)
	}
}

/*
Accept Messages once a connection is established. Dedicated routine for each connection
Arguments: TCP Socket Connection Object
Returns/Exits on Node Failure Only
*/
func handleRequests(conn net.Conn) {

	r := bufio.NewReader(conn)
	for {
		fmt.Println("Server waiting for message")
		data, err := r.ReadString('\n')
		if !checkError(err, "handleRequests") {
			conn.Close()
			totalNodes--
			return
		}

		// Reset Timer if received message was a heartbeat
		go updateHBFlag(data, timer)
		fmt.Println("Received:", data)
		go processRequest(data)
	}
}

/*
Initializes timer for heartbeat timeout
*/
func heartbeatChecker() {

	// Waits for commmand from main routine till all connections/nodes are ready. (On load only)
	fmt.Println("Heartbeat channel", <-chanStartHBCheck)
	fmt.Println("heartbeat checker starts")
	heartbeatCheckerElection()
}

func heartbeatCheckerElection() {
	fmt.Println("heartbeat checker starts")
	timeOut := heartbeatTimeOut + r.Intn(2*heartbeatTimeOut)
	fmt.Println("timeout:", timeOut)
	fmt.Println("time Now:", time.Now())
	timer = time.NewTimer(time.Duration(timeOut) * time.Millisecond)
	<-timer.C
	fmt.Println("-------------------------------------")
	fmt.Println("NO HEARTBEAT RECIEVED WITHIN TIMEOUT")
	fmt.Println("-------------------------------------")
	fmt.Println("time Now:", time.Now())
	go initiateElection(nodeID)
	heartbeatCheckerElection()
}

/*
Scans all incoming messages to check if it is Heartbeat from leader
*/
func updateHBFlag(data string, timer *time.Timer) {
	dataSlice := strings.Fields(data)
	if strings.Compare(dataSlice[0], "ThisIsHeartbeat") == 0 {
		resetTimer()
		s := getState()
		t, _ := strconv.Atoi(dataSlice[2])
		fmt.Println("leader:", dataSlice[1])
		insertTableState(t, s.votedFor, dataSlice[1], 0)
	}
}

/*
Reset the timer when a heartbeat is received
*/
func resetTimer() {
	fmt.Println("Reset Timeout")
	timeOut := heartbeatTimeOut + r.Intn(2*heartbeatTimeOut)
	fmt.Println(timeOut)
	fmt.Println("time Now:", time.Now())
	timer.Stop()
	timer.Reset(time.Duration(timeOut) * time.Millisecond)
}

/*
Kills the timer once node wins the election
*/
func killTimer() {
	fmt.Println("Kill Timer")
	timer.Stop()
}

/*
Dummy method for testing only
*/
func testingConn(c net.Conn) {
	fmt.Println("testing conn start")
	data := "An Acknowledgment"
	_, err := c.Write([]byte(data))
	checkError(err, "testingConn")
	fmt.Println("testing conn end")
}
