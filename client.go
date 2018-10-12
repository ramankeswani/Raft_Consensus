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
func client(port int, myPort int, connChan chan connection, remoteNodeID string) {

	if port == 0 {
		fmt.Fprintf(os.Stderr, "Port not given")
		os.Exit(1)
	}

	address := "127.0.0.1:" + strconv.Itoa(port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	checkError(err, "client")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	c := connection{
		nodeID: remoteNodeID,
		conn:   conn,
	}
	connChan <- c
	fmt.Println("inserted into channel")
	checkError(err, "client")
	_, err = conn.Write([]byte("First request" + strconv.Itoa(myPort) + "\n"))
	checkError(err, "client")

}

/*
Sends the message string to other Node's Server
Arguments: message to be passed, and remotenodeID
Invoked for all outgoing messages
*/
func sendMessage(message string, remoteNodeID string) {
	fmt.Println("Client sendMessage Starts:", remoteNodeID)
	c := connMap[remoteNodeID]
	_, err := c.conn.Write([]byte(message))
	checkError(err, "sendMessage")
	fmt.Println("Client sendMessage Ends")
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
