package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

var myPort int

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Provide IP, usage %s port ", os.Args[0])
		os.Exit(1)
	}

	myPort, _ = strconv.Atoi(os.Args[1])
	cS := make(chan int)
	go server(myPort, cS)

	tableCluster()
	ns := getNodesFromDB()
	fmt.Println("Recieved nodes from DB as: ")
	fmt.Printf("%+v", ns)

	time.Sleep(10 * time.Second)
	go sendConnectionRequest(ns)

	fmt.Println(<-cS)
	fmt.Println("Program ends")

}

func sendConnectionRequest(ns nodes) {

	fmt.Println("\n sendConnectionRequest Starts")
	for n := range ns {
		if myPort != ns[n].port {
			fmt.Println(ns[n])
			time.Sleep(3 * time.Second)
			go client(ns[n].port)
		}
	}

	fmt.Println("\n sendConnectionRequest Ends")
}

func checkError(err error, src string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s Fatal error: %s \n", src, err.Error())
		os.Exit(1)
	}
}
