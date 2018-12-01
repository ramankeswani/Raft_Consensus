cd
cd Raft_Consensus
sudo rm -rf databases
sudo mkdir databases
sudo git pull
sudo go build
sudo nohup ./Raft_Consensus $1 $2 0 2>&1 > console_$2.txt &
sudo nohup ./Raft_Consensus $3 $4 0 2>&1 > console_$4.txt &
sudo nohup ./Raft_Consensus $5 $6 0 2>&1 > console_$6.txt &