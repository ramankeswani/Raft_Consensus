package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

var connections []*net.TCPConn
var timeout = 20

func heartbeat(otherNodes nodes, myNodeID string) {

	time.Sleep(3 * time.Second)
	fmt.Println("Heartbeat Starts")
	establishConn(otherNodes)
	for {
		time.Sleep(time.Duration(timeout) * time.Second)
		for conn := range connections {
			go sendHeartbeat(connections[conn])
		}
		fmt.Println("Heartbeat Ends")
	}
}

func sendHeartbeat(conn *net.TCPConn) {
	_, err := conn.Write([]byte("ThisIsHeartbeat"))
	checkError(err, "Heartbeat")
}

func establishConn(ns nodes) {

	for n := range ns {
		address := "127.0.0.1:" + strconv.Itoa(ns[n].port)
		tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
		checkError(err, "Heartbeat")

		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		checkError(err, "Heartbeat")

		_, err = conn.Write([]byte("First from heartbeat"))

		connections = append(connections, conn)
	}
}
