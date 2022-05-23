package main

import (
	"encoding/json"
	"fmt"
	"my/AlienDB/region/sqlite"
	"net"

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
	quitStatement     = 3
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
		case quitStatement:
			QuitChan <- "quit"
		}

		if err != nil {
			res.Error = err.Error()
		}
		res.Error = "success"
		res.Data = msg

		if err != nil {
			QuitChan <- err.Error()
		} else {
			output <- res

			select {
			case rec := <-input:
				//fmt.Println("rec", rec)
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
				//fmt.Println("sqlExec", msg, err)
				if err != nil {
					res.Error = err.Error()
				}
				res.Error = ""
				res.Data = msg
				fmt.Println(">>> 得到结果: ", res)
				if err != nil {
					QuitChan <- "err"
				} else {
					output <- res
				}

			}
		}
	}
}

func input(connMaster net.Conn, input chan receive) {
	//通信得到结果
	for {
		msg := make([]byte, 255)
		msgRead, err := connMaster.Read(msg)
		data := make([]byte, msgRead)
		copy(data, msg)
		if msgRead == 0 || err != nil {
			panic(err)
		} else {
			request := make(map[string]string)
			json.Unmarshal(data, &request)
			fmt.Println(">>> 收到请求: ", request)
			if request["kind"] == "select" {
				temp := &receive{
					sqlStatement: request["sql"],
					sqlType:      queryStatement,
				}
				input <- *temp
			} else {
				temp := &receive{
					sqlStatement: request["sql"],
					sqlType:      nonQueryStatement,
				}
				input <- *temp
			}
		}
	}
}

func output(connMaster net.Conn, output chan result) {
	for {
		select {
		case outPutMsg := <-output:
			//通信返回结果
			fmt.Println(">>> 返回给 Master: ", outPutMsg)
			msgStr, _ := json.Marshal(outPutMsg)
			//fmt.Println("msgStr", msgStr)
			_, err := connMaster.Write(msgStr)
			if err != nil {
				panic(err)
			}
		}

	}
}

func main() {
	defer sqlite.Close()
	fmt.Println(">>> Region 启动中...")
	//var endpoints = []string{"localhost:2222"}
	StatementChannel := make(chan receive, 500)
	OutputChannel := make(chan result, 500)
	QuitChan = make(chan string)
	//ser, err := sqlite.NewServiceRegister(endpoints, "/db/region_01", "localhost:8000", 6, 5)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//go ser.ListenLeaseRespChan()

	connMaster := sqlite.ConnectToMaster()
	defer connMaster.Close()
	go input(connMaster, StatementChannel)
	go handle(StatementChannel, OutputChannel)
	go output(connMaster, OutputChannel)
	for {

		toQuit := <-QuitChan
		if toQuit == "quit" {
			fmt.Println("bye")
		} else {
			fmt.Println(toQuit)
			fmt.Println("unexpected error , terminated")

		}

	}
}
