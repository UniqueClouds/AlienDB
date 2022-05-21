package sqlite

import (
	"database/sql"
	"os"
)

var (
	db *sql.DB
)

func init() {
	var err error
	os.Create("./foo2.db")
	db, err = sql.Open("sqlite3", "./foo2.db")
	checkErr(err)
}
