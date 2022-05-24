package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
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
	dropAllTable()
	checkErr(err)
}

func dropAllTable() {
	rows, err := Query("SELECT name FROM sqlite_master WHERE type = 'table';")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(rows)
	var names []string
	for _, val := range rows {
		//fmt.Println("val", val["name"])
		names = append(names, val["name"].(string))
	}
	for _, name := range names {
		if strings.Compare(name, "sqlite_sequence") == 0 {
			continue
		}
		fmt.Println("drop table ", name)
		_, err := db.Exec("DROP table if exists " + name + ";")
		//fmt.Println("exec", exec)
		if err != nil {
			fmt.Println("err: ", err)
		}
	}
}
