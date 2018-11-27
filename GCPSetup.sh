cd
sudo wget https://storage.googleapis.com/golang/go1.10.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.10.linux-amd64.tar.gz
sudo export PATH=$PATH:/usr/local/go/bin:${HOME}/go/bin
cd go/src
git clone https://github.com/mattn/go-sqlite3.git
go build github.com/mattn/go-sqlite3
cd
cd Raft_Consensus
go build