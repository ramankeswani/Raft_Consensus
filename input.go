package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/papertrail/go-tail/follower"
)

func userInput(connMap map[string]connection) {
	fmt.Println("User Input Starts")
	for key, value := range connMap {
		fmt.Println("User Input key: " + key + " value: " + value.nodeID)
	}

	t, _ := follower.New("userinput.txt", follower.Config{
		Whence: io.SeekEnd,
		Offset: 0,
		Reopen: true,
	})

	//reader := bufio.NewReader(os.Stdin)
	//for {

	for line := range t.Lines() {
		//text, _ := reader.ReadString('\n')
		text := line.String()
		fmt.Println("user entered: ", text)
		s := getState()
		if strings.Compare(myNodeID, s.leader) != 0 {
			fmt.Println("leader", s.leader)
			go contactLeader(connMap[s.leader], text)
		} else {
			go appendEntryInit(strings.TrimRight(text, "\n"))
		}
	}
}

func contactLeader(conn connection, text string) {
	fmt.Println("Contact Leader Starts")
	message := myNodeID + " " + AppendEntryFromClient + " " + text + "\n"
	_, err := conn.conn.Write([]byte(message))
	checkError(err, "contactLeader")
	fmt.Println("Contact Leader Ends")
}
