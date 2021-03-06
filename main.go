package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var logTag = "*"
var myPort string
var nodeID string
var otherNodes nodes
var connChan chan connection
var chanStartHBCheck chan string
var chanMessage chan string
var chanLatestLog chan string
var chanSyncResp chan string
var chanAppendResp chan string
var isRecovering bool
var mutexConnMap = &sync.Mutex{}
var mutexChanMap = &sync.Mutex{}

type connection struct {
	nodeID string
	//conn   *net.TCPConn
	conn   net.Conn
	status bool
}

var totalNodes int
var connMap map[string]connection
var chanMap map[string]chan string

/*
Entry Point for the Application
Usage: ./main port nodeID
*/
func main() {

	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Provide IP, usage %s port nodeID RecoverFlag", os.Args[0])
		os.Exit(1)
	}
	myPort = os.Args[1]
	nodeID = os.Args[2]
	initLog(nodeID)
	connChan = make(chan connection)
	chanStartHBCheck = make(chan string)
	connMap = make(map[string]connection)
	chanMap = make(map[string]chan string)
	chanSyncResp = make(chan string)
	chanLatestLog = make(chan string)
	chanAppendResp = make(chan string, 500)
	tableCluster(nodeID)
	ns := getNodesFromDB()
	totalNodes = getTotalNodes()

	fmt.Println("I AM: ", nodeID+" my port"+os.Args[1])
	fmt.Println("------------------------------")
	c := make(chan int)
	setFirstStartIndex()

	// Starting Server
	go server(myPort, nodeID)
	if strings.Compare(os.Args[3], "0") == 0 {
		time.Sleep(45 * time.Second)
	} else {
		time.Sleep(10 * time.Second)
	}
	populateOtherNodes(ns)
	go sendConnectionRequest(otherNodes)
	for range otherNodes {
		conn := <-connChan
		fmt.Println("connection channel received status: " + strconv.FormatBool(conn.status))
		if conn.status {
			mutexConnMap.Lock()
			connMap[conn.nodeID] = conn
			mutexConnMap.Unlock()
			chanTemp := make(chan string)
			mutexChanMap.Lock()
			chanMap[conn.nodeID] = chanTemp
			mutexChanMap.Unlock()
			go sendMessage(conn.nodeID)
		}
	}

	// TO-DO: Assuming node ALPHA to be the leader.
	/* if strings.Compare(nodeID, "ALPHA") == 0 {
		s := state{
			currentTerm: 1,
		}
		go heartbeat(otherNodes, nodeID, connMap, s)
	} */

	// Check if Node is Recovering
	/* 	if strings.Compare(os.Args[3], "1") == 0 {
		initRecovery(myPort)
	} */

	s := getState()
	logFile("recover", "state main: "+strconv.Itoa(s.startIndex))
	if s.startIndex > 1 {
		initRecovery(myPort)
		isRecovering = true
	}

	go userInput(connMap)
	chanStartHBCheck <- "start"
	//isCommited = make(map[int]bool)

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
		go client(ns[n].port, myPort, connChan, ns[n].nodeID, ns[n].address)
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
