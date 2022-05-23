package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var (
	db      *sql.DB
	logFile *os.File
)

func logInit() {
	var err error
	logFile, err = os.OpenFile("./c.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Panic("打开日志文件异常")
	}
	log.SetOutput(logFile)
}

func PrintLog(sql string) {
	log.Println(sql)
}

func init() {
	logInit()
	var err error
	db, err = sql.Open("sqlite3", "./foo2.db")
	fmt.Println(db)
	checkErr(err)
}
