package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

/*
Creates TCP Socket connection objects with other nodes in the cluster on load
Arguments: Remote Port, Server Port, Connection Channel (To send back Socket Object to Main Routine), remote NodeID,
Channel to Initialize HashMap of connection Objects
Returns when Node Crashes Only
*/
var remoteID string

func client(port int, myPort string, connChan chan connection, remoteNodeID string, remoteAddress string) {

	if port == 0 {
		fmt.Fprintf(os.Stderr, "Port not given")
		os.Exit(1)
	}

	remoteID = remoteNodeID
	address := remoteAddress + ":" + strconv.Itoa(port)
	//	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	//checkError(err, "client")
	//conn, err := net.DialTCP("tcp", nil, tcpAddr)
	fmt.Println("inside client address: " + address)
	conn, err := net.Dial("tcp", address)
	st := checkError(err, "client after dial")
	c := connection{
		nodeID: remoteNodeID,
		conn:   conn,
		status: st,
	}
	connChan <- c
	fmt.Println("inserted into channel")
	if st {
		_, err = conn.Write([]byte("First request" + myPort + "\n"))
		checkError(err, "client")
	}

}

/*
Sends the message string to all other Nodes
Arguments: message to be passed
Invoked for all broadcast outgoing messages
*/
func sendMessageToAll(message string) {
	fmt.Println("Send Message To All Starts Map Len:", len(connMap))
	for _, conn := range connMap {
		_, err := conn.conn.Write([]byte(message))
		checkError(err, "Send Message To All")
	}
	fmt.Println("Send Message To All Ends")
}

/* func waitForMessages(conn *net.TCPConn) {
	for {
		r := bufio.NewReader(conn)
		data, err := r.ReadString('\n')
		checkError(err, "waitForMessages")
		fmt.Println("received at client: ", data)
		go processRequest(data, conn)
	}
}
*/
