package main

import (
	"database/sql"
	"fmt"
	"os"

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
}

type nodes []node

var insertStateStmt *sql.Stmt
var delStateStmt *sql.Stmt
var dbName string

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
	/* _, err = insertStatement.Exec("NODE4", "127.0.0.1", 7001)
	_, err = insertStatement.Exec("NODE5", "127.0.0.1", 7002) */

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

	createStatement, err := db.Prepare("CREATE TABLE state (currentTerm integer, votedFor text, leader text)")
	createStatement.Exec()
	checkErr(err)

	insertStateStmt, err = db.Prepare("INSERT INTO state(currentTerm, votedFor, leader) values(?,?,?)")
	insertStateStmt.Exec(0, "", "")

	delStateStmt, err = db.Prepare("DELETE FROM state")
	checkErr(err)
}

func insertTableState(currentTerm int, votedFor string, leader string) (res bool) {

	delStateStmt.Exec()
	_, err := insertStateStmt.Exec(currentTerm, votedFor, leader)
	checkErr(err)

	if err != nil {
		return false
	}
	return true
}

func getState() state {
	var currTerm int
	var votedFor string
	var leader string
	s := state{
		currentTerm: -1,
		votedFor:    "",
		leader:      "",
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
		err = row.Scan(&currTerm, &votedFor, &leader)
	}
	if err != nil {
		return s
	}
	s = state{
		currentTerm: currTerm,
		votedFor:    votedFor,
		leader:      leader,
	}
	row.Close()
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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
