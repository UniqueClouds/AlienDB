package sqlite

import (
	"database/sql"
	"fmt"
)

var (
	db *sql.DB
)

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./foo2.db")
	fmt.Println(db)
	checkErr(err)
}
