package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func client(port int, myPort int) {

	if port == 0 {
		fmt.Fprintf(os.Stderr, "Port not given")
		os.Exit(1)
	}

	address := "127.0.0.1:" + strconv.Itoa(port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	checkError(err, "client")

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err, "client")

	_, err = conn.Write([]byte("First request" + strconv.Itoa(myPort)))
	checkError(err, "client")

	request := make([]byte, 128)

	result, err := conn.Read(request)
	checkError(err, "client")

	fmt.Println(string(request[:result]))

}
