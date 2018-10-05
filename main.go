package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var myPort int
var nodeID string
var otherNodes nodes
var connChan chan connection
var chanStartHBCheck chan string

type connection struct {
	nodeID string
	conn   *net.TCPConn
}

var connMap map[string]connection

var chanConnMap chan map[string]connection

// Entry Point for the Application
// Usage go run port nodeID
func main() {

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Provide IP, usage %s port nodeID ", os.Args[0])
		os.Exit(1)
	}

	myPort, _ = strconv.Atoi(os.Args[1])
	nodeID = os.Args[2]
	connChan = make(chan connection)
	chanStartHBCheck = make(chan string)
	connMap := make(map[string]connection)
	chanConnMap = make(chan map[string]connection)

	tableCluster(nodeID)
	fmt.Println("I AM: ", nodeID)
	fmt.Println("------------------------------")

	c := make(chan int)
	// Starting Server
	go server(myPort, c, nodeID)

	ns := getNodesFromDB()
	populateOtherNodes(ns)

	// Send Intitial Connection Requests to Other Nodes
	go sendConnectionRequest(otherNodes)

	for range otherNodes {
		conn := <-connChan
		fmt.Println("connection channel received", conn)
		connMap[conn.nodeID] = conn
	}

	// TO-DO: Assuming node ALPHA to be the leader.
	if strings.Compare(nodeID, "ALPHA") == 0 {
		go heartbeat(otherNodes, nodeID, connMap)
	}
	chanConnMap <- connMap
	chanStartHBCheck <- "start"
	// Blocking to keep the main routine alive forever
	fmt.Println(<-c)
	fmt.Println("Program ends")

}

func sendConnectionRequest(ns nodes) {

	// Sleep till other nodes are ready
	time.Sleep(3 * time.Second)

	fmt.Println("\nSendConnectionRequest Starts")
	for n := range ns {
		fmt.Println(ns[n])
		go client(ns[n].port, myPort, connChan, ns[n].nodeID, chanConnMap)
	}
	fmt.Println("SendConnectionRequest Ends")
	fmt.Println("------------------------------")

}

func checkError(err error, src string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s Fatal error: %s \n", src, err.Error())
		os.Exit(1)
	}
}

func populateOtherNodes(ns nodes) {
	for n := range ns {
		if nodeID != ns[n].nodeID {
			otherNodes = append(otherNodes, ns[n])
		}
	}
}
