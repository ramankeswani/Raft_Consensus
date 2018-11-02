package main

import (
	"database/sql"
	"fmt"
	"os"
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
var dbName string
var mutex = &sync.Mutex{}
var mutexLogTable = &sync.Mutex{}
var logIndex int

func tableCluster(nodeID string) {

	dbName = "./databases/" + nodeID + ".db"
	os.Remove(dbName)

	db, err := sql.Open("sqlite3", dbName)
	checkErr(err)

	createStatement, err := db.Prepare("CREATE TABLE cluster (nodeID text PRIMARY KEY, address text, port integer)")
	createStatement.Exec()
	checkErr(err)

	insertStatement, err := db.Prepare("INSERT INTO cluster(nodeID, address, port) values(?,?,?)")
	checkErr(err)

	_, err = insertStatement.Exec("ALPHA", "127.0.0.1", 5000)
	_, err = insertStatement.Exec("BETA", "127.0.0.1", 6000)
	_, err = insertStatement.Exec("GAMMA", "127.0.0.1", 7000)
	_, err = insertStatement.Exec("NODE4", "127.0.0.1", 7001)
	_, err = insertStatement.Exec("NODE5", "127.0.0.1", 7002)

	checkErr(err)

	fmt.Println("-----------------------------------------")

	tableState()
}

func tableState() {

	/* 	os.Remove("./raft.db")

	   	db, err := sql.Open("sqlite3", "./raft.db")
		   checkErr(err) */
	fmt.Println("dbname:", dbName)
	db, err := sql.Open("sqlite3", dbName)

	createStatement, err := db.Prepare("CREATE TABLE state (currentTerm integer, votedFor text, leader text, commitIndex integer)")
	createStatement.Exec()
	checkErr(err)

	insertStateStmt, err = db.Prepare("INSERT INTO state(currentTerm, votedFor, leader, commitIndex) values(?,?,?,?)")
	insertStateStmt.Exec(0, "", "", 0)

	delStateStmt, err = db.Prepare("DELETE FROM state")
	checkErr(err)
	tableLog()
}

func insertTableState(currentTerm int, votedFor string, leader string, commitIndex int) (res bool) {
	mutex.Lock()
	delStateStmt.Exec()
	_, err := insertStateStmt.Exec(currentTerm, votedFor, leader, commitIndex)
	checkErr(err)
	mutex.Unlock()
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
	s := state{
		currentTerm: -1,
		votedFor:    "",
		leader:      "",
		commitIndex: -1,
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
		err = row.Scan(&currTerm, &votedFor, &leader, &commIndex)
	}
	if err != nil {
		return s
	}
	s = state{
		currentTerm: currTerm,
		votedFor:    votedFor,
		leader:      leader,
		commitIndex: commIndex,
	}
	row.Close()
	mutex.Unlock()
	return s
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
	return ns

}

func tableLog() {
	db, err := sql.Open("sqlite3", dbName)

	createStatementLog, err := db.Prepare("CREATE TABLE log (logIndex integer PRIMARY KEY AUTOINCREMENT, term integer, command text, votes integer, commited integer)")
	createStatementLog.Exec()
	checkErr(err)

	insertLogStmt, err = db.Prepare("INSERT INTO log(term, command, votes, commited) values(?,?,?,?)")
	checkErr(err)
	updateLogStmt, err = db.Prepare("UPDATE log SET votes=? where logIndex=?")
	checkErr(err)
	commitLogStmt, err = db.Prepare("UPDATE log SET commited=? where logIndex=?")
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func insertLogTable(term int, command string, votes int) (prevLogIndex int, prevLogTerm int) {
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
	row, err := db.Query("SELECT * from log order by logIndex desc limit 1")
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
