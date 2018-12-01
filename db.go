package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type node struct {
	nodeID  string
	address string
	port    int
}

type state struct {
	currentTerm int
	votedFor    string
	leader      string
	commitIndex int
	startIndex  int
}

type log struct {
	logIndex   int
	term       int
	command    string
	votes      int
	isCommited int
}

type nodes []node

var insertStateStmt *sql.Stmt
var delStateStmt *sql.Stmt
var insertLogStmt *sql.Stmt
var updateLogStmt *sql.Stmt
var commitLogStmt *sql.Stmt
var createNextIndexStmt *sql.Stmt
var insertNextIndexStmt *sql.Stmt
var updateNextIndexStmt *sql.Stmt
var dbName string
var mutex = &sync.Mutex{}
var mutexLogTable = &sync.Mutex{}
var mutexLogTableInsert = &sync.Mutex{}
var logIndex int

func tableCluster(nodeID string) {

	fmt.Println("table cluster: " + nodeID)
	dbName = "./databases/" + nodeID + ".db"
	//os.Remove(dbName)

	db, err := sql.Open("sqlite3", dbName)
	checkErr(err)

	dropStmt, err := db.Prepare("drop table if exists cluster")
	dropStmt.Exec()

	createStatement, err := db.Prepare("CREATE TABLE cluster (nodeID text PRIMARY KEY, address text, port integer)")
	createStatement.Exec()

	insertStatement, err := db.Prepare("INSERT INTO cluster(nodeID, address, port) values(?,?,?)")
	checkErr(err)
	// for i, j := 1, 5001; i <= 8; i, j = i+1, j+1 {
	// 	_, err = insertStatement.Exec("node"+strconv.Itoa(i), "127.0.0.1", j)
	// }

	var ipMap = make(map[string]string)
	ipMap["node1"] = "10.142.0.4"
	ipMap["node2"] = "10.142.0.5"
	ipMap["node3"] = "10.142.0.6"
	ipMap["node4"] = "10.142.0.7"
	ipMap["node5"] = "10.142.0.8"
	ipMap["node6"] = "10.142.0.9"
	ipMap["node7"] = "10.142.0.10"
	ipMap["node8"] = "10.142.0.12"
	ipMap["node17"] = "10.150.0.3"
	ipMap["node18"] = "10.150.0.4"
	ipMap["node19"] = "10.150.0.5"
	ipMap["node20"] = "10.150.0.6"
	ipMap["node21"] = "10.150.0.7"
	ipMap["node22"] = "10.150.0.8"
	ipMap["node23"] = "10.150.0.9"
	ipMap["node24"] = "10.150.0.2"
	ipMap["node33"] = "10.158.0.2"
	ipMap["node34"] = "10.158.0.3"
	ipMap["node35"] = "10.158.0.4"
	ipMap["node36"] = "10.158.0.5"
	ipMap["node37"] = "10.158.0.6"
	ipMap["node38"] = "10.158.0.7"
	ipMap["node39"] = "10.158.0.8"
	ipMap["node40"] = "10.158.0.9"

	ind := 1
	j := 5000

	for _, v := range ipMap {
		for i := 0; i < 3; i++ {
			name := "node" + strconv.Itoa(i*24+ind)
			_, err = insertStatement.Exec(name, v, j+i*24+ind)
		}
		ind++
	}

	checkErr(err)

	fmt.Println("-----------------------------------------")

	tableState()
}

func getIP(nodeM string) (ip string) {
	fmt.Println("get my ip starts: " + dbName + " SELECT * FROM cluster where nodeID = " + nodeM)
	db, err := sql.Open("sqlite3", dbName)
	checkErr(err)
	row, err := db.Query("SELECT address FROM cluster where cluster.nodeID = \"" + nodeM + "\"")
	checkErr(err)
	for row.Next() {
		err = row.Scan(&ip)
		fmt.Println("ip is" + ip)
	}
	db.Close()
	fmt.Println("get my ip ends: " + nodeM)
	return ip
}

func tableState() {

	/* 	os.Remove("./raft.db")

	   	db, err := sql.Open("sqlite3", "./raft.db")
		   checkErr(err) */
	fmt.Println("dbname:", dbName)
	db, err := sql.Open("sqlite3", dbName)

	createStatement, err := db.Prepare("CREATE TABLE IF NOT EXISTS state (currentTerm integer, votedFor text, leader text, commitIndex integer, startIndex integer)")
	createStatement.Exec()
	checkErr(err)

	insertStateStmt, err = db.Prepare("INSERT INTO state(currentTerm, votedFor, leader, commitIndex, startIndex) values(?,?,?,?,?)")
	//insertStateStmt.Exec(0, "", "", 0, 0)

	delStateStmt, err = db.Prepare("DELETE FROM state")
	checkErr(err)
	tableLog()
	//tableNextIndex()
}

func setFirstStartIndex() {
	insertTableState(0, "", "", 0, 1)
}

func insertTableState(currentTerm int, votedFor string, leader string, commitIndex int, startIndex int) (res bool) {
	startIndex = getState().startIndex + startIndex
	delStateStmt.Exec()
	_, err := insertStateStmt.Exec(currentTerm, votedFor, leader, commitIndex, startIndex)
	checkErr(err)
	if err != nil {
		return false
	}
	return true
}

func getState() state {
	mutex.Lock()
	var currTerm int
	var votedFor string
	var leader string
	var commIndex int
	var startIndex int
	s := state{
		currentTerm: -1,
		votedFor:    "",
		leader:      "",
		commitIndex: -1,
		startIndex:  0,
	}
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return s
	}
	row, err := db.Query("SELECT * FROM state order by currentTerm DESC limit 1")
	if err != nil {
		return s
	}
	for row.Next() {
		err = row.Scan(&currTerm, &votedFor, &leader, &commIndex, &startIndex)
	}
	if err != nil {
		return s
	}
	s = state{
		currentTerm: currTerm,
		votedFor:    votedFor,
		leader:      leader,
		commitIndex: commIndex,
		startIndex:  startIndex,
	}
	row.Close()
	db.Close()
	mutex.Unlock()
	return s
}

func getTotalNodes() int {
	fmt.Println("get Total Nodes start")
	var count int
	db, err := sql.Open("sqlite3", dbName)
	checkErr(err)
	rows, err := db.Query("SELECT count(*) FROM cluster")
	checkErr(err)
	for rows.Next() {
		err = rows.Scan(&count)
		fmt.Println("count:" + strconv.Itoa(count))
		checkErr(err)
	}
	fmt.Println("get Total Nodes end")
	db.Close()
	return count

}

func getNodesFromDB() nodes {

	fmt.Println("getNodesFromDB Starts:", dbName)

	ns := nodes{}
	var nodeID string
	var address string
	var port int

	db, err := sql.Open("sqlite3", dbName)
	checkErr(err)

	rows, err := db.Query("SELECT * FROM cluster")
	checkErr(err)

	for rows.Next() {
		err = rows.Scan(&nodeID, &address, &port)
		fmt.Println(nodeID, address, port)
		n := node{
			nodeID:  nodeID,
			address: address,
			port:    port,
		}
		ns = append(ns, n)
	}

	fmt.Println("getNodesFromDB Ends")
	fmt.Println("------------------------------")
	db.Close()
	return ns

}

func tableLog() {
	db, err := sql.Open("sqlite3", dbName)

	createStatementLog, err := db.Prepare("CREATE TABLE IF NOT EXISTS log (logIndex integer PRIMARY KEY AUTOINCREMENT, term integer, command text, votes integer, commited integer)")
	createStatementLog.Exec()
	checkErr(err)

	insertLogStmt, err = db.Prepare("INSERT INTO log(term, command, votes, commited) values(?,?,?,?)")
	checkErr(err)
	updateLogStmt, err = db.Prepare("UPDATE log SET votes=? where logIndex=?")
	checkErr(err)
	commitLogStmt, err = db.Prepare("UPDATE log SET commited=? where logIndex=?")
	checkErr(err)
}

func tableNextIndex() {
	db, err := sql.Open("sqlite3", dbName)

	createNextIndexStmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS nextIndex (nodeId text, index integer)")
	createNextIndexStmt.Exec()
	checkErr(err)

	insertNextIndexStmt, err = db.Prepare("INSERT INTO nextIndex(nodeId, index) values (?,?)")
	insertNextIndexStmt.Exec()
	checkErr(err)

	updateNextIndexStmt, err = db.Prepare("UPDATE nextIndex SET index = ? where nodeId = ?")
	checkErr(err)

	for node := range otherNodes {
		insertNextIndexStmt.Exec(otherNodes[node].nodeID, 1)
	}
}

func updateNextIndex(nodeID string, index int) (res bool) {
	logFile("commit", "updateNextIndex starts nodeID: "+nodeID+" index: "+strconv.Itoa(index)+"\n")
	rows, err := updateNextIndexStmt.Exec(nodeID, index)
	checkErr(err)
	lastID, _ := rows.LastInsertId()
	logFile("commit", "updateNextIndex lastID: "+strconv.Itoa(int(lastID))+"\n")
	if lastID <= 0 {
		return false
	}
	logFile("commit", "updateNextIndex ends\n")
	return true
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func insertLogTable(term int, command string, votes int) (prevLogIndex int, prevLogTerm int) {
	mutexLogTableInsert.Lock()
	logFile("append", "insertLogTable starts command: "+command+"\n")
	res, err := insertLogStmt.Exec(term, command, votes, 0)
	temp, _ := res.RowsAffected()
	logFile("append", "insertLogTable err null is null: "+strconv.FormatBool(err == nil)+" res rows affected: "+strconv.Itoa(int(temp))+"\n")
	checkErr(err)
	logIndex, _ := res.LastInsertId()
	logCurrent := getLogTable(int(logIndex))
	logFile("append", "insertLogTable logcurrent  index: "+strconv.Itoa(logCurrent.logIndex)+" command: "+logCurrent.command+"\n")
	if err != nil {
		return -1, -1
	}
	logIndex, _ = res.LastInsertId()
	logPrevious := getLogTable(int(logCurrent.logIndex - 1))
	logFile("append", "insertLogTable previous tuple command: "+logPrevious.command+" index: "+strconv.Itoa(logPrevious.logIndex)+"\n")
	mutexLogTableInsert.Unlock()
	return int(logPrevious.logIndex), logPrevious.term
}

func incrementVoteCount(index int) (count int) {
	logFile("append", "increment vote count index: "+strconv.Itoa(index)+"\n")
	mutexLogTable.Lock()
	logFile("append", "increment lock started\n")
	db, err := sql.Open("sqlite3", dbName)
	checkErr(err)
	row, err := db.Query("select * from log where logIndex = " + strconv.Itoa(index))
	l := getLogTable(index)
	logFile("append", "vote: "+strconv.Itoa(l.votes)+" term: "+strconv.Itoa(l.term)+" command: "+l.command+" logIndex: "+strconv.Itoa(logIndex)+"\n")
	status := updateLogTable(l.votes+1, index)
	logFile("append", "increment lock ended\n")
	row.Close()
	db.Close()
	mutexLogTable.Unlock()
	if status {
		return l.votes + 1
	}
	return -1
}

func commitLog(index int) (result bool) {
	logFile("commit", "commitLog Starts index: "+strconv.Itoa(index)+"\n")
	res, err := commitLogStmt.Exec(1, index)
	checkErr(err)
	ra, _ := res.RowsAffected()
	logFile("commit", "commitLog rows affected: "+strconv.Itoa(int(ra))+"\n")
	logFile("commit", "commitLog Ends index: "+strconv.Itoa(index)+"\n")
	if ra == 1 {
		return true
	}
	return false
}

func updateLogTable(index int, votes int) (result bool) {
	logFile("append", "updateLogTable Starts index: "+strconv.Itoa(index)+" votes: "+strconv.Itoa(votes)+"\n")
	rows, err := updateLogStmt.Exec(index, votes)
	lastID, _ := rows.LastInsertId()
	rowsAff, _ := rows.RowsAffected()
	logFile("append", "updateLogTable lastID: "+strconv.Itoa(int(lastID))+" rowsAff: "+strconv.Itoa(int(rowsAff))+"\n")
	checkErr(err)
	logFile("append", "updateLogTable calling getLogTable\n")
	_ = getLogTable(int(lastID))
	return true
}

func getLogTable(logIndex int) (l log) {
	logFile("append", "getLogTable starts logindex: "+strconv.Itoa(logIndex)+"\n")
	db, err := sql.Open("sqlite3", dbName)
	rows, err := db.Query("SELECT * FROM log where logIndex = " + strconv.Itoa(logIndex))
	//rows, err := db.Query("SELECT * FROM log where logIndex = 1")
	var lIndex int
	var t int
	var command string
	var v int
	var c int
	l = log{
		logIndex: -1,
	}
	for rows.Next() {
		err = rows.Scan(&lIndex, &t, &command, &v, &c)
		checkErr(err)
		logFile("append", "getLogTable Inside Loop: lIndex: "+strconv.Itoa(lIndex)+" command: "+command+" votes: "+strconv.Itoa(v)+" isCommited: "+strconv.Itoa(c)+"\n")
		l = log{
			logIndex:   lIndex,
			term:       t,
			command:    command,
			votes:      v,
			isCommited: c,
		}
	}
	rows.Close()
	db.Close()
	logFile("append", "getLogTable ends\n")
	return l
}

func getLatestLog() (l log) {
	logFile("commit", "getLatestLog() starts\n")
	db, err := sql.Open("sqlite3", dbName)
	checkErr(err)
	row, err := db.Query("SELECT logIndex from log order by logIndex desc limit 1")
	var id int
	for row.Next() {
		err = row.Scan(&id)
		fmt.Println("commit", "getLatestLog() index: "+strconv.Itoa(id)+"\n")
	}
	l = getLogTable(id)
	db.Close()
	logFile("commit", "getLatestLog() ends\n")
	return l
}
