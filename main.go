package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Provide IP, usage %s port ", os.Args[0])
		os.Exit(1)
	}

	clustertable()

	c := make(chan int)
	port, _ := strconv.Atoi(os.Args[1])
	go server(port, c)
	time.Sleep(3 * time.Second)

	go client(port, c)

	fmt.Println(<-c)
	fmt.Println(<-c)

}

func checkError(err error, src string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s Fatal error: %s \n", src, err.Error())
		os.Exit(1)
	}
}
