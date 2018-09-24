pkill konsole
go build main.go server.go client.go heartbeat.go db.go

konsole --noclose -e "./main 5000 ALPHA" & 
sleep 1
konsole --noclose -e "./main 6000 BETA" &
sleep 1
konsole --noclose -e "./main 7000 GAMMA" &