package sqlite

import (
	_ "github.com/mattn/go-sqlite3"
)

func Exec(sqlString string) (string, error) {
	//db, err := sql.Open("sqlite3", "./foo2.db")
	//defer db.Close()
	//if err != nil {
	//	return "database open failed", err
	//}
	_, err := db.Exec(sqlString) // ignore_security_alert

	if err != nil {
		return "database open fail check if there is a database on the region", err
	}
	return "success", nil
}
