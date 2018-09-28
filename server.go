package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// Server Module for Node
func server(myPort int, c chan int, nodeID string) {

	service := ":" + strconv.Itoa(myPort)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err, "server")

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err, "server")

	timer := time.NewTimer(time.Duration(heartbeatTimeOut) * time.Second)
	// Routine to track HeartBeat Timeout
	go heartbeatChecker(timer)

	for {
		request := make([]byte, 128)

		fmt.Println("Server Ready to Accept")
		conn, err := listener.Accept()
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
		go handleRequests(conn, timer)

	}

}

// Accept Messages once a connection is established
func handleRequests(conn net.Conn, timer *time.Timer) {

	for {
		request := make([]byte, 128)
		inp, err := conn.Read(request)
		checkError(err, "Client HandleRequests")
		data := string(request[:inp])
		// Reset Timer if received message was a heartbeat
		go updateHBFlag(data, timer)
		fmt.Println("Received:", data)
	}
}

func heartbeatChecker(timer *time.Timer) {

	<-timer.C
	fmt.Println("------------------------------")
	fmt.Println("NO HEARTBEAT RECIEVED WITHIN TIMEOUT")
	fmt.Println("------------------------------")
}

func updateHBFlag(data string, timer *time.Timer) {
	if strings.Compare(data, "ThisIsHeartbeat") == 0 {
		timer.Reset(10 * time.Second)
	}
}
