package main

/*
Sends the message string to other Node's Server
Arguments: message to be passed, and remotenodeID
Invoked for all outgoing messages
*/
func sendMessage(remoteID string) {
	for {
		//fmt.Println("Client sendMessage Starts remoteID:", remoteID)
		mutexChanMap.Lock()
		ch := chanMap[remoteID]
		mutexChanMap.Unlock()
		m := <-ch
		//fmt.Println("received on channel at client: ", m)
		mutexConnMap.Lock()
		c := connMap[remoteID]
		mutexConnMap.Unlock()
		_, err := c.conn.Write([]byte(m))
		checkError(err, "sendMessage")
		//fmt.Println("Client sendMessage Ends")
	}
}
