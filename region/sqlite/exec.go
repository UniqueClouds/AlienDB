package sqlite

import (
	_ "github.com/mattn/go-sqlite3"
)

func Exec(sqlString string) ([]map[string]interface{}, error) {
	//db, err := sql.Open("sqlite3", "./foo2.db")
	//defer db.Close()
	//if err != nil {
	//	return "database open failed", err
	//}
	_, err := db.Exec(sqlString) // ignore_security_alert

	if err != nil {
		return nil, err
	}
	return nil, nil
}
