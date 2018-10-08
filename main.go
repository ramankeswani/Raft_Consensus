package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
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

var totalNodes int
var connMap map[string]connection
var chanConnMap chan map[string]connection

/*
Entry Point for the Application
Usage: go run port nodeID
*/
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
	go server(myPort, nodeID, chanConnMap)

	ns := getNodesFromDB()
	populateOtherNodes(ns)
	go sendConnectionRequest(otherNodes)
	for range otherNodes {
		conn := <-connChan
		fmt.Println("connection channel received", conn)
		connMap[conn.nodeID] = conn
	}

	// TO-DO: Assuming node ALPHA to be the leader.
	/* if strings.Compare(nodeID, "ALPHA") == 0 {
		s := state{
			currentTerm: 1,
		}
		go heartbeat(otherNodes, nodeID, connMap, s)
	} */

	// Once for server, once for server helper and once for client
	chanConnMap <- connMap
	chanConnMap <- connMap
	chanConnMap <- connMap
	totalNodes = len(connMap) + 1
	go userInput(connMap)
	chanStartHBCheck <- "start"

	// Blocking to keep the main routine alive forever
	fmt.Println(<-c)
	fmt.Println("Program ends")

}

/*
Create TCP sockets with all other nodes in the cluster
Arguements: Node Objects
Returns when connection objects are established
*/
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

func checkError(err error, src string) (status bool) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s Fatal error: %s \n", src, err.Error())
		return false
	}
	return true
}

func populateOtherNodes(ns nodes) {
	for n := range ns {
		if nodeID != ns[n].nodeID {
			otherNodes = append(otherNodes, ns[n])
		}
	}
}
