package sqlite

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestExec(t *testing.T) {

	//Exec("CREATE TABLE `userinfo` (`till_id` INTEGER PRIMARY KEY AUTOINCREMENT, `client_id` VARCHAR(64) NULL, `first_name` VARCHAR(255) NOT NULL, `last_name` VARCHAR(255) NOT NULL, `guid` VARCHAR(255) NULL, `dob` DATETIME NULL, `type` VARCHAR(1))")
	//rows, err := Query("SELECT * FROM userinfo")
	//fmt.Println(rows, err)
	Exec("CREATE TABLE `userinfo2` (`till_id` INTEGER PRIMARY KEY AUTOINCREMENT, `client_id` VARCHAR(64) NULL, `first_name` VARCHAR(255) NOT NULL, `last_name` VARCHAR(255) NOT NULL, `guid` VARCHAR(255) NULL, `dob` DATETIME NULL, `type` VARCHAR(1));INSERT INTO userinfo2(client_id, first_name, last_name) values(1,'2','2');INSERT INTO userinfo(client_id, first_name, last_name) values(1,'3','3');", "userinfo")
	rows, err := Query("SELECT name FROM sqlite_master WHERE type = 'table';")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rows)
	var names []string
	for _, val := range rows {
		fmt.Println("val", val["name"])
		names = append(names, val["name"].(string))
	}
	for _, name := range names {
		if strings.Compare(name, "sqlite_sequence") == 0 {
			continue
		}
		fmt.Println("drop table ", name)
		exec, err := db.Exec("DROP table if exists " + name + ";")
		fmt.Println("exec", exec)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println(Query("select name from sqlite_master where type = 'table';"))
}

func TestBack(t *testing.T) {
	ctx := context.Background()
	BackupSQLite3(ctx)
}

func TestRestore(t *testing.T) {
	Restore("./backup/userinfo.txt")
}
