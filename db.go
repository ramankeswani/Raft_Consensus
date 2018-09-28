package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type node struct {
	nodeID  string
	address string
	port    int
}

type nodes []node

func tableCluster() {

	/* os.Remove("./raft.db")

	db, err := sql.Open("sqlite3", "./raft.db")
	checkErr(err)

	createStatement, err := db.Prepare("CREATE TABLE cluster (nodeID text PRIMARY KEY, address text, port integer)")
	createStatement.Exec()
	checkErr(err)

	insertStatement, err := db.Prepare("INSERT INTO cluster(nodeID, address, port) values(?,?,?)")
	checkErr(err)

	_, err = insertStatement.Exec("ALPHA", "127.0.0.1", 5000)
	_, err = insertStatement.Exec("BETA", "127.0.0.1", 6000)
	_, err = insertStatement.Exec("GAMMA", "127.0.0.1", 7000)
	checkErr(err) */
}

func getNodesFromDB() nodes {

	fmt.Println("getNodesFromDB Starts")

	ns := nodes{}
	var nodeID string
	var address string
	var port int

	db, err := sql.Open("sqlite3", "./raft.db")
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
