package main

import (
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
var totalNodes = len(connections) + 1

// Server Module for Node
func server(myPort int, c chan int, nodeID string) {

	source = rand.NewSource(time.Now().UnixNano())
	r = rand.New(source)

	service := ":" + strconv.Itoa(myPort)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err, "server")

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err, "server")

	// Routine to track HeartBeat Timeout
	go heartbeatChecker()

	for {
		request := make([]byte, 5120)

		fmt.Println("Server Ready to Accept")
		conn, err := listener.Accept()
		//defer conn.Close()
		if err != nil {
			fmt.Println("Error at server while accepting")
			continue
		}
		fmt.Println("Server accepted From: ", conn.RemoteAddr())

		inp, err := conn.Read(request)
		checkError(err, "client")

		fmt.Println("Server Received (Just After Accept):", string(request[:inp]))

		daytime := time.Now().String()
		conn.Write([]byte(daytime + " nodeID is " + nodeID))
		//conn.Close()

		// Routine to accept messages
		go handleRequests(conn)

		conn.Write([]byte("-----testing 1------"))
		conn.Write([]byte("-----testing 2------"))
		conn.Write([]byte("-----testing 3------"))

	}

}

// Accept Messages once a connection is established
func handleRequests(conn net.Conn) {

	for {
		//defer conn.Close()
		request := make([]byte, 5120)
		fmt.Println("Server waiting for message")
		inp, err := conn.Read(request)
		checkError(err, "Server HandleRequests")
		data := string(request[:inp])
		// Reset Timer if received message was a heartbeat
		go updateHBFlag(data, timer)
		fmt.Println("Received:", data)

		conn.Write([]byte("ack: " + data))
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
	fmt.Println("------------------------------")
	fmt.Println("NO HEARTBEAT RECIEVED WITHIN TIMEOUT")
	fmt.Println("------------------------------")
	go initiateElection(nodeID)
}

func updateHBFlag(data string, timer *time.Timer) {
	if strings.Compare(data, "ThisIsHeartbeat") == 0 {
		timeOut := heartbeatTimeOut + r.Intn(3*heartbeatTimeOut) + r.Intn(3*heartbeatTimeOut)
		fmt.Println(timeOut)
		timer.Reset(time.Duration(timeOut) * time.Millisecond)
	}
}
