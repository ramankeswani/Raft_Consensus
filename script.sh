pkill konsole
go build main.go server.go serverHelper.go client.go heartbeat.go db.go RaftConstants.go

konsole --noclose -e "./main 5000 ALPHA" & 
echo $!
konsole --noclose -e "./main 6000 BETA" &
konsole --noclose -e "./main 7000 GAMMA" &