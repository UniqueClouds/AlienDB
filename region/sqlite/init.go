package sqlite

import (
	"database/sql"
)

var (
	db *sql.DB
)

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./foo3.db")
	//fmt.Println(db)
	checkErr(err)
}
