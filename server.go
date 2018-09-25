package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func server(myPort int, c chan int, nodeID string) {

	service := ":" + strconv.Itoa(myPort)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err, "server")

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err, "server")

	timer := time.NewTimer(15 * time.Second)
	go heartbeatChecker(timer)

	for {
		request := make([]byte, 128)
		fmt.Println("Server Ready to Accept")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error at server while accepting")
			continue
		}
		fmt.Println("Server accepted From: ", conn.RemoteAddr())

		inp, err := conn.Read(request)
		checkError(err, "client")

		fmt.Println(string(request[:inp]))

		daytime := time.Now().String()
		conn.Write([]byte(daytime + " nodeID is " + nodeID))
		//conn.Close()

		go handleRequests(conn, timer)

	}

}

func handleRequests(conn net.Conn, timer *time.Timer) {

	for {
		request := make([]byte, 128)
		inp, err := conn.Read(request)
		checkError(err, "Client HandleRequests")
		data := string(request[:inp])
		go updateHBFlag(data, timer)
		fmt.Println(data)
	}
}

func heartbeatChecker(timer *time.Timer) {

	<-timer.C
	fmt.Println("NO HEARTBEAT SORRY")
}

func updateHBFlag(data string, timer *time.Timer) {
	if strings.Compare(data, "ThisIsHeartbeat") == 0 {
		timer.Reset(10 * time.Second)
	}
}
