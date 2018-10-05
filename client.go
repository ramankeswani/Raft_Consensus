package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
)

var cMap map[string]connection

// Client module
func client(port int, myPort int, connChan chan connection, remoteNodeID string, chanConnMap chan map[string]connection) {

	if port == 0 {
		fmt.Fprintf(os.Stderr, "Port not given")
		os.Exit(1)
	}
	go initMap(chanConnMap)
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

	//go waitForMessages(conn)

	//request := make([]byte, 128)
	/*
		r := bufio.NewReader(conn)
		data, err := r.ReadString('\n')
		checkError(err, "client")

		fmt.Println("Handshake Response: ", data)
		 for {
			r := bufio.NewReader(conn)
			data, err := r.ReadString('\n')
			checkError(err, "client")

			fmt.Println("Handshake Response: ", data)
		} */

}

/*
Sends the message string to other Node's Server
Arguments message to be passed, and TCP socket connection
Invoked for all outgoing messages
*/
func sendMessage(message string, remoteNodeID string) {
	fmt.Println("Client | sendMessage Starts |", remoteNodeID)
	c := cMap[remoteNodeID]
	_, err := c.conn.Write([]byte(message))
	checkError(err, "sendMessage")
	fmt.Println("Client | sendMessage Ends")
}

func initMap(chanConnMap chan map[string]connection) {
	cMap = <-chanConnMap
}

func waitForMessages(conn *net.TCPConn) {
	for {
		r := bufio.NewReader(conn)
		data, err := r.ReadString('\n')
		checkError(err, "waitForMessages")
		fmt.Println("received at client: ", data)
		go processRequest(data, conn)
	}
}
