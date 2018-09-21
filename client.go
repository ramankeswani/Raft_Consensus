package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func client(port int, c chan int) {

	if port == 0 {
		fmt.Fprintf(os.Stderr, "Port not given")
		os.Exit(1)
	}

	address := "127.0.0.1:" + strconv.Itoa(port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	checkError(err, "client")

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err, "client")

	time.Sleep(2 * time.Second)
	_, err = conn.Write([]byte("Client sending data"))
	checkError(err, "client")

	request := make([]byte, 128)

	result, err := conn.Read(request)
	checkError(err, "client")

	fmt.Println(string(request[:result]))

	c <- 10
	os.Exit(0)
}
