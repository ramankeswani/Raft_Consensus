package main

import (
	"fmt"
	"net"
	"time"
)

// Initialise Heartbeat
func heartbeat(otherNodes nodes, myNodeID string, connections []*net.TCPConn) {

	//time.Sleep(3 * time.Second)

	fmt.Println("Heartbeat Starts")

	for {
		break
		time.Sleep(time.Duration(heartbeatInterval) * time.Millisecond)
		for conn := range connections {
			// Send heartbeat after interval time
			go sendHeartbeat(connections[conn])
		}
		fmt.Println("Heartbeat Ends")
	}
}

func sendHeartbeat(conn *net.TCPConn) {
	fmt.Println("sendHeartbeat")
	_, err := conn.Write([]byte("ThisIsHeartbeat"))
	checkError(err, "Heartbeat")
}

/* func establishConn(ns nodes) {

	for n := range ns {
		address := "127.0.0.1:" + strconv.Itoa(ns[n].port)
		tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
		checkError(err, "Heartbeat")

		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		checkError(err, "Heartbeat")

		_, err = conn.Write([]byte("First Message from heartbeat"))
		connections = append(connections, conn)
	}
} */
