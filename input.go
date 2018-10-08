package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func userInput(connMap map[string]connection) {
	for key, value := range connMap {
		fmt.Println("key: " + key + " value: " + value.nodeID)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		fmt.Println("user entered: ", text)
		s := getState()
		if strings.Compare(myNodeID, s.leader) != 0 {
			fmt.Println("leader", s.leader)
			go contactLeader(connMap[s.leader], text)
		}
	}
}

func contactLeader(conn connection, text string) {
	fmt.Println("Contact Leader Starts")
	message := myNodeID + " " + AppendEntryFromClient + " " + text
	_, err := conn.conn.Write([]byte(message))
	checkError(err, "contactLeader")
	fmt.Println("Contact Leader Ends")
}
