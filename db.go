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

	createStatement, err := db.Prepare("CREATE TABLE state (currentTerm integer, votedFor text)")
	createStatement.Exec()
	checkErr(err)

	insertStateStmt, err = db.Prepare("INSERT INTO state(currentTerm, votedFor) values(?,?)")
	insertStateStmt.Exec(0, "")

	delStateStmt, err = db.Prepare("DELETE FROM state")
	checkErr(err)
}

func insertTableState(currentTerm int, votedFor string) (res bool) {

	delStateStmt.Exec()
	_, err := insertStateStmt.Exec(currentTerm, votedFor)
	checkErr(err)

	if err != nil {
		return false
	}
	return true
}

func getState() (currentTerm int, votedFor string) {
	var currTerm int
	var voteFor string
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return -1, ""
	}
	row, err := db.Query("SELECT * FROM state order by currentTerm DESC limit 1")
	if err != nil {
		return -2, ""
	}
	for row.Next() {
		err = row.Scan(&currTerm, &voteFor)
	}
	if err != nil {
		return -3, ""
	}
	row.Close()
	return currTerm, voteFor
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
