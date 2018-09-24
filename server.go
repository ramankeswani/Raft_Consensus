package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func server(myPort int, c chan int, nodeID string) {

	service := ":" + strconv.Itoa(myPort)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err, "server")

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err, "server")

	for {
		request := make([]byte, 128)
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

		go handleRequests(conn)
	}
}

func handleRequests(conn net.Conn) {

	for {
		request := make([]byte, 128)
		inp, err := conn.Read(request)
		checkError(err, "Client HandleRequests")
		fmt.Println(string(request[:inp]))
	}
}
