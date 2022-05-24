package sqlite

import (
	"context"
	"testing"
)

func TestExec(t *testing.T) {

	//Exec("CREATE TABLE `userinfo` (`till_id` INTEGER PRIMARY KEY AUTOINCREMENT, `client_id` VARCHAR(64) NULL, `first_name` VARCHAR(255) NOT NULL, `last_name` VARCHAR(255) NOT NULL, `guid` VARCHAR(255) NULL, `dob` DATETIME NULL, `type` VARCHAR(1))")
	//rows, err := Query("SELECT * FROM userinfo")
	//fmt.Println(rows, err)
	Exec("CREATE TABLE `userinfo2` (`till_id` INTEGER PRIMARY KEY AUTOINCREMENT, `client_id` VARCHAR(64) NULL, `first_name` VARCHAR(255) NOT NULL, `last_name` VARCHAR(255) NOT NULL, `guid` VARCHAR(255) NULL, `dob` DATETIME NULL, `type` VARCHAR(1));INSERT INTO userinfo2(client_id, first_name, last_name) values(1,'2','2');INSERT INTO userinfo(client_id, first_name, last_name) values(1,'3','3');", "userinfo")

}

func TestBack(t *testing.T) {
	ctx := context.Background()
	BackupSQLite3(ctx)
}

func TestRestore(t *testing.T) {
	Restore("./backup/userinfo.txt")
}
