pkill konsole
go build main.go server.go serverHelper.go client.go recoverNode.go log.go heartbeat.go db.go RaftConstants.go input.go appendEntryRPC.go sendMessage.go

konsole --noclose -e "./main 5000 ALPHA 0" & 
echo $!
konsole --noclose -e "./main 6000 BETA 0" &
konsole --noclose -e "./main 7000 GAMMA 0" &
konsole --noclose -e "./main 7001 NODE4 0" &
konsole --noclose -e "./main 7002 NODE5 0" &