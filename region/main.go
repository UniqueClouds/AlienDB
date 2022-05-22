package main

import (
	"encoding/json"
	"fmt"
	"my/AlienDB/region/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

var QuitChan chan string

type receive struct {
	sqlStatement string
	sqlType      int
}

type result struct {
	Error string                   `json:"error"`
	Data  []map[string]interface{} `json:"data"`
}

const (
	queryStatement    = 1
	nonQueryStatement = 2
)

func handle(input chan receive, output chan result) {
	for {
		rec := <-input
		var (
			res result
			msg []map[string]interface{}
			err error
		)

		switch rec.sqlType {
		case queryStatement:
			msg, err = sqlite.Query(rec.sqlStatement)

		case nonQueryStatement:
			msg, err = sqlite.Exec(rec.sqlStatement)
		}

		if err != nil {
			res.Error = err.Error()
		}
		res.Error = "success"
		res.Data = msg

		if err != nil {
			QuitChan <- "err"
		} else {
			output <- res
		}
	}
}

func input(input chan receive) {
	//通信得到结果
	temp := &receive{
		//sqlStatement: "INSERT INTO userinfo(client_id, first_name, last_name) values(4,'2','3')",
		//sqlStatement: "CREATE TABLE `userinfo` (`till_id` INTEGER PRIMARY KEY AUTOINCREMENT, `client_id` VARCHAR(64) NULL, `first_name` VARCHAR(255) NOT NULL, `last_name` VARCHAR(255) NOT NULL, `guid` VARCHAR(255) NULL, `dob` DATETIME NULL, `type` VARCHAR(1))",
		sqlStatement: "SELECT * FROM userinfo",
		sqlType:      queryStatement,
	}
	input <- *temp
}

func output(output chan result) {
	outPutMsg := <-output
	//通信返回结果
	test, err := json.Marshal(outPutMsg)
	if err != nil {
		fmt.Println(test)
	}
	var res result
	err = json.Unmarshal(test, &res)
	if err != nil {
		fmt.Println(res)
	}

	fmt.Println(outPutMsg)
}

func main() {
	defer sqlite.Close()
	StatementChannel := make(chan receive, 500)
	OutputChannel := make(chan result, 500)
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
