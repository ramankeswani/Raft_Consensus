cd
sudo wget https://storage.googleapis.com/golang/go1.10.linux-amd64.tar.gz
sudo tar -xzf go1.10.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin:${HOME}/go/bin
cd go/src
sudo mkdir github.com
cd github.com
sudo mkdir mattn
cd mattn
git clone https://github.com/mattn/go-sqlite3.git
go build github.com/mattn/go-sqlite3
cd
cd Raft_Consensus
go build
sudo mkdir databases
sudo mkdir logs