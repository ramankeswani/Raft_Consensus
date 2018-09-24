package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func server(myPort int, c chan int) {

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
		fmt.Println("server accepted", conn.RemoteAddr())

		inp, err := conn.Read(request)
		checkError(err, "client")

		fmt.Println(string(request[:inp]))

		daytime := time.Now().String()
		conn.Write([]byte(daytime + " port is " + strconv.Itoa(myPort)))
		conn.Close()
	}
}
