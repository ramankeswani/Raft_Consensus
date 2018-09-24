package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var myPort int
var nodeID string
var otherNodes nodes

func main() {

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Provide IP, usage %s port nodeID ", os.Args[0])
		os.Exit(1)
	}

	myPort, _ = strconv.Atoi(os.Args[1])
	nodeID = os.Args[2]
	c := make(chan int)
	go server(myPort, c, nodeID)

	tableCluster()
	ns := getNodesFromDB()
	populateOtherNodes(ns)
	fmt.Println("Data From DB:")
	fmt.Printf("%+v\n", ns)
	fmt.Printf("compareing")
	fmt.Println(strings.Compare(nodeID, "ALPHA") == 0)
	time.Sleep(10 * time.Second)
	go sendConnectionRequest(otherNodes)
	if strings.Compare(nodeID, "ALPHA") == 0 {
		go heartbeat(otherNodes, nodeID)
	}
	fmt.Println(<-c)
	fmt.Println("Program ends")

}

func sendConnectionRequest(ns nodes) {

	fmt.Println("\n sendConnectionRequest Starts")
	for n := range ns {
		time.Sleep(3 * time.Second)
		fmt.Println(ns[n])
		go client(ns[n].port)
	}
	fmt.Println("\n sendConnectionRequest Ends")
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
