package main

import (
	"fmt"
	"net"
	"time"
)

func server(myPort int, c chan int) {

	service := ":5000"

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

		//inp, err := ioutil.ReadAll(conn)
		inp, err := conn.Read(request)
		checkError(err, "client")

		fmt.Println(string(request[:inp]))

		daytime := time.Now().String()
		conn.Write([]byte(daytime))
		//conn.Close()
	}
}
