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
	logIndex int
	term     int
	command  string
	votes    int
}

type nodes []node

var insertStateStmt *sql.Stmt
var delStateStmt *sql.Stmt
var insertLogStmt *sql.Stmt
var updateLogStmt *sql.Stmt
var dbName string
var mutex = &sync.Mutex{}

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

	createStatementLog, err := db.Prepare("CREATE TABLE log (logIndex integer PRIMARY KEY AUTOINCREMENT, term integer, command text, votes integer)")
	createStatementLog.Exec()
	checkErr(err)

	insertLogStmt, err = db.Prepare("INSERT INTO log(term, command, votes) values(?,?,?)")
	checkErr(err)
	updateLogStmt, err = db.Prepare("UPDATE log SET votes=? where logIndex=?")
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func insertLogTable(term int, command string, votes int) (prevLogIndex int, prevLogTerm int) {
	res, err := insertLogStmt.Exec(term, command, votes)
	checkErr(err)

	if err != nil {
		return -1, -1
	}
	id, _ := res.LastInsertId()
	log := getLogTable(int(id) - 1)
	return int(id) - 1, log.term
}

func updateLogTable(index int, votes int) (result bool) {
	_, err := updateLogStmt.Exec(index, votes)
	checkErr(err)

	if err != nil {
		return false
	}
	return true
}

func getLogTable(logIndex int) (l log) {
	db, err := sql.Open("sqlite3", dbName)
	rows, err := db.Query("SELECT * FROM log where logIndex = " + strconv.Itoa(logIndex))
	var lIndex int
	var t int
	var command string
	var v int
	l = log{
		logIndex: -1,
	}
	for rows.Next() {
		err = rows.Scan(&lIndex, &t, &command, &t, &v)
		checkErr(err)
		l = log{
			logIndex: lIndex,
			term:     t,
			command:  command,
			votes:    v,
		}
	}
	return l
}
