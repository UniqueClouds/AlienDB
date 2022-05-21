package main

import (
	"fmt"
	"my/AlienDB/region/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

var QuitChan chan string

type receive struct {
	sqlStatement string
	sqlType      int
}

const (
	queryStatement    = 1
	nonQueryStatement = 2
)

func handle(input chan receive, output chan string) {
	for {
		rec := <-input
		var (
			msg string
			err error
		)

		switch rec.sqlType {
		case queryStatement:
			msg, err = sqlite.Query(rec.sqlStatement)
		case nonQueryStatement:
			msg, err = sqlite.Exec(rec.sqlStatement)
		}
		if err != nil {
			QuitChan <- "err"
		} else {
			output <- msg
		}
	}
}

func input(input chan receive) {
	temp := &receive{
		//sqlStatement: "INSERT INTO userinfo(client_id, first_name, last_name) values(1,'1','1')",
		//sqlStatement: "CREATE TABLE `userinfo` (`till_id` INTEGER PRIMARY KEY AUTOINCREMENT, `client_id` VARCHAR(64) NULL, `first_name` VARCHAR(255) NOT NULL, `last_name` VARCHAR(255) NOT NULL, `guid` VARCHAR(255) NULL, `dob` DATETIME NULL, `type` VARCHAR(1))",
		sqlStatement: "SELECT * FROM userinfo",
		sqlType:      queryStatement,
	}
	input <- *temp
}

func output(output chan string) {
	outPutMsg := <-output
	fmt.Println(outPutMsg)
}

func main() {
	defer sqlite.Close()
	StatementChannel := make(chan receive, 500)
	OutputChannel := make(chan string, 500)
	QuitChan = make(chan string)

	// SOCKET msg input StatementChannel
	go input(StatementChannel)
	go handle(StatementChannel, OutputChannel)
	// SOCKET msg output OutputChannel
	go output(OutputChannel)

	toQuit := <-QuitChan
	if toQuit == "err" {
		fmt.Println("unexpected err occur")
	} else {
		fmt.Println("bye")
	}

}
