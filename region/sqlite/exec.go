package sqlite

import (
	"database/sql"
	"os"
)

func Exec(sqlString string) (string, error) {
	os.Create("./foo2.db")
	db, err := sql.Open("sqlite3", "./foo2.db")
	if err != nil {
		return "database open fail check if there is a database on the region", err
	}

	_, err = db.Exec(sqlString) // ignore_security_alert

	if err != nil {
		return "database open fail check if there is a database on the region", err
	}
	return "success", nil
}
