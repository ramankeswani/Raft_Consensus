package main

/*
Sends the message string to other Node's Server
Arguments: message to be passed, and remotenodeID
Invoked for all outgoing messages
*/
func sendMessage(remoteID string) {
	for {
		//fmt.Println("Client sendMessage Starts remoteID:", remoteID)
		m := <-chanMap[remoteID]
		//fmt.Println("received on channel at client: ", m)
		c := connMap[remoteID]
		_, err := c.conn.Write([]byte(m))
		checkError(err, "sendMessage")
		//fmt.Println("Client sendMessage Ends")
	}
}
