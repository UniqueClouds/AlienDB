package sqlite

import (
	"database/sql"
	//"fmt"
)

var (
	db *sql.DB
)

//
//const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
//
//func RandStringBytes(n int) string {
//	rand.Seed(time.Now().UnixNano())
//	b := make([]byte, n)
//	for i := range b {
//		b[i] = letterBytes[rand.Intn(len(letterBytes))]
//	}
//	return string(b)
//}
func init() {
	var err error
	//dbName := RandStringBytes(4) + ".db"
	//fmt.Println(">>> dbName: ", dbName)
	db, err = sql.Open("sqlite3", "./foo2.db")
	db.Exec("delete from sqlite_master where type in ('table', 'index', 'trigger');")
	//fmt.Println(db)
	checkErr(err)
}
