package main

import (
	"fmt"
	"strconv"
	"time"
)

/*
Invoked only by leader to initialize and send heartbeat messages
Arguments: Other Nodes in Cluster, Self NodeID, hashmap of connection objects
Returns when node crashes or new leader is elected
*/
var i int

func heartbeat(otherNodes nodes, myNodeID string, connMap map[string]connection, s state) {
	i++
	fmt.Println("Heartbeat Starts i:", i)
	for leader {
		// For testing heartbeat failure only
		/* if i == 1 && strings.Compare("ALPHA", myNodeID) == 0 {
			break
		} */
		if !leader {
			return
		}
		time.Sleep(time.Duration(heartbeatInterval) * time.Millisecond)
		go sendHeartbeat(myNodeID, s.currentTerm)
		fmt.Println("Heartbeat Ends")
	}
}

/*
Sends out actual heartbeat message
Arguments: TCP Socket Connection Object, Self/Leader NodeID
Returns ASAP after sending heartbeat
*/
func sendHeartbeat(myNodeID string, term int) {
	message := "ThisIsHeartbeat" + " " + myNodeID + " " + strconv.Itoa(term) + "\n"
	for _, c := range chanMap {
		c <- message
	}
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
