package main

import (
	"fmt"
	"net"
	"time"
)

/*
Invoked only by leader to initialize and send heartbeat messages
Arguments: Other Nodes in Cluster, Self NodeID, hashmap of connection objects
Returns when node crashes or new leader is elected
*/
func heartbeat(otherNodes nodes, myNodeID string, connMap map[string]connection) {

	fmt.Println("Heartbeat Starts")
	for {
		// For testing heartbeat failure only
		break
		time.Sleep(time.Duration(heartbeatInterval) * time.Millisecond)
		for _, conn := range connMap {
			go sendHeartbeat(conn.conn, myNodeID)
		}
		fmt.Println("Heartbeat Ends")
	}
}

/*
Sends out actual heartbeat message
Arguments: TCP Socket Connection Object, Self/Leader NodeID
Returns ASAP after sending heartbeat
*/
func sendHeartbeat(conn *net.TCPConn, myNodeID string) {
	fmt.Println("sendHeartbeat")
	_, err := conn.Write([]byte("ThisIsHeartbeat" + " " + myNodeID + "\n"))
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
