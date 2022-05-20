package sqlite

import "testing"

func TestExec(t *testing.T) {
	Exec("CREATE TABLE `userinfo` (`till_id` INTEGER PRIMARY KEY AUTOINCREMENT, `client_id` VARCHAR(64) NULL, `first_name` VARCHAR(255) NOT NULL, `last_name` VARCHAR(255) NOT NULL, `guid` VARCHAR(255) NULL, `dob` DATETIME NULL, `type` VARCHAR(1))")
}
