package sqlite

import (
	"database/sql"
	"fmt"
	"log"
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
	row, err := db.Query("SELECT name FROM sqlite_master WHERE type = 'table';")
	//fmt.Println(res)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	for row.Next() {
		var name string
		row.Scan(&name)
		fmt.Println("drop Table Name ", name)
		db.Exec("DROP TABLE IF EXISTS " + name)
	}
	//fmt.Println(db)
	checkErr(err)
}
