package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func clustertable() {

	db, err := sql.Open("sqlite3", "./foo.db")
	checkErr(err)

	statement, err := db.Prepare("CREATE TABLE cluster (nodeid text PRIMARY KEY, address text, port integer)")
	statement.Exec()

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
