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
var votes = 0
var totalNodes int
var myNodeID string

// Server Module for Node
func server(myPort int, c chan int, nodeID string) {

	totalNodes = len(connMap) + 1
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
		//request := make([]byte, 1000000)

		fmt.Println("Server Ready to Accept")
		conn, err := listener.Accept()
		conn.Write([]byte("ack from inside\n"))
		//defer conn.Close()
		if err != nil {
			fmt.Println("Error at server while accepting")
			continue
		}
		fmt.Println("Server accepted From: ", conn.RemoteAddr())

		/* inp, err := conn.Read(request)
		checkError(err, "client") */
		/* r := bufio.NewReader(conn)
		data, err := r.ReadString('\n')
		checkError(err, "server")

		fmt.Println("Server Received (Just After Accept):", data) */

		/* daytime := time.Now().String()
		conn.Write([]byte(daytime + " nodeID is " + nodeID + "\n"))
		//conn.Close()

		conn.Write([]byte("-----testing 1------\n"))
		conn.Write([]byte("-----testing 2------\n"))
		conn.Write([]byte("-----testing 3------\n")) */

		// Routine to accept messages

		go handleRequests(conn)

	}

}

// Accept Messages once a connection is established
func handleRequests(conn net.Conn) {

	r := bufio.NewReader(conn)
	/* 	data, err := r.ReadString('\n')
	   	checkError(err, "server")

	   	fmt.Println("Server Received (Just After Accept):", data)
	*/
	for {
		//defer conn.Close()

		/* request := make([]byte, 1000000)
		fmt.Println("Server waiting for message")
		inp, err := conn.Read(request)
		checkError(err, "Server HandleRequests")
		fmt.Println("inp:", inp)
		data := string(request[:inp]) */

		fmt.Println("Server waiting for message")
		data, err := r.ReadString('\n')
		checkError(err, "handleRequests")

		// Reset Timer if received message was a heartbeat
		go updateHBFlag(data, timer)
		fmt.Println("Received:", data)
		//testingConn(conn)

		conn.Write([]byte("ack\n"))
		go processRequest(data, conn)
	}
}

func heartbeatChecker() {

	fmt.Println(<-chanStartHBCheck)
	fmt.Println("heartbeat checker starts")
	timeOut := heartbeatTimeOut + r.Intn(3*heartbeatTimeOut) + r.Intn(3*heartbeatTimeOut)
	fmt.Println(timeOut)
	timer = time.NewTimer(time.Duration(timeOut) * time.Millisecond)
	<-timer.C
	fmt.Println("-------------------------------------")
	fmt.Println("NO HEARTBEAT RECIEVED WITHIN TIMEOUT")
	fmt.Println("-------------------------------------")
	go initiateElection(nodeID)
}

func updateHBFlag(data string, timer *time.Timer) {
	dataSlice := strings.Fields(data)
	if strings.Compare(dataSlice[0], "ThisIsHeartbeat") == 0 {
		timeOut := heartbeatTimeOut + r.Intn(3*heartbeatTimeOut) + r.Intn(3*heartbeatTimeOut)
		fmt.Println(timeOut)
		timer.Reset(time.Duration(timeOut) * time.Millisecond)
	}
}

func resetTimer() {
	timeOut := heartbeatTimeOut + r.Intn(3*heartbeatTimeOut) + r.Intn(3*heartbeatTimeOut)
	fmt.Println(timeOut)
	timer.Reset(time.Duration(timeOut) * time.Millisecond)
}

func testingConn(c net.Conn) {
	fmt.Println("testing conn start")
	data := "An Acknowledgment"
	_, err := c.Write([]byte(data))
	checkError(err, "testingConn")
	fmt.Println("testing conn end")
}
