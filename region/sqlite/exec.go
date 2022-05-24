package sqlite

import (
	"bufio"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var backFile *os.File

func logInit(s string) {
	var err error
	os.MkdirAll("./backup", 0755)
	backFile, err = os.OpenFile("./backup/"+s+".txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Panic("打开日志文件异常")
	}
}

func Exec(sqlString string, tableName string) ([]map[string]interface{}, error) {
	//db, err := sql.Open("sqlite3", "./foo2.db")
	//defer db.Close()
	//if err != nil {
	//	return "database open failed", err
	//}
	defer backFile.Close()

	logInit(tableName)
	_, err := db.Exec(sqlString) // ignore_security_alert
	write := bufio.NewWriter(backFile)
	write.WriteString(sqlString)
	write.Flush()
	if err != nil {
		return nil, err
	}

	return nil, nil
}
