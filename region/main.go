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

func main() {

	StatementChannel := make(chan receive, 500)
	OutputChannel := make(chan string, 500)
	QuitChan = make(chan string)

	// SOCKET msg input StatementChannel
	go handle(StatementChannel, OutputChannel)
	// SOCKET msg output OutputChannel

	toQuit := <-QuitChan
	if toQuit == "err" {
		fmt.Println("unexpected err occur")
	} else {
		fmt.Println("bye")
	}

}
