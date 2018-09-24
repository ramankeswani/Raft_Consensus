package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func clustertable() {
	os.Remove("./raft.db")
	db, err := sql.Open("sqlite3", "./raft.db")
	checkErr(err)

	createStatement, err := db.Prepare("CREATE TABLE cluster (nodeid text PRIMARY KEY, address text, port integer)")
	createStatement.Exec()
	checkErr(err)

	insertStatement, err := db.Prepare("INSERT INTO cluster(nodeid, address, port) values(?,?,?)")
	checkErr(err)

	res, err := insertStatement.Exec("ALPHA", "127.0.0.1", 5000)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)
	fmt.Println("Last inserted id:", id)

	updateStatment, err := db.Prepare("UPDATE CLUSTER SET nodeid=? where nodeid=?")
	checkErr(err)

	res, err = updateStatment.Exec("BETA", "ALPHA")
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println("Rows affected:", affect)

	rows, err := db.Query("SELECT * FROM cluster")
	checkErr(err)

	var nodeid string
	var address string
	var port int

	for rows.Next() {
		err = rows.Scan(&nodeid, &address, &port)
		fmt.Println(nodeid, address, port)
	}

	rows.Close()

	db.Close()
}

/*
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
*/
