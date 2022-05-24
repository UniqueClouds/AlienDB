package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"my/AlienDB/region/sqlite"
	"net"
	"os"
)

var QuitChan chan string

type receive struct {
	sqlStatement string
	sqlType      int
	tableName    string
	ipAddress    string
	file         []string // 用于 copy table
}

type regionRequest struct {
	TableName string
	IpAddress string
	Kind      string
	Sql       string
	File      []string // 用于 copy table
}
type result struct {
	Error     string                   `json:"error"`
	Data      []map[string]interface{} `json:"data"`
	TableList []string                 `json:"tableList"`
	Message   string                   `json:"message"`
	ClientIP  string                   `json:"clientIP"`
	File      []string                 `json:"file"`
}

const (
	queryStatement    = 1
	nonQueryStatement = 2
	quitStatement     = 3
	copyStatement     = 4
	joinStatement     = 5
)

func main() {
	defer sqlite.Close()
	StatementChannel := make(chan receive, 500)
	OutputChannel := make(chan result, 500)
	QuitChan = make(chan string)
	fmt.Println(">>> Region 启动中...")
	connMaster := sqlite.ConnectToMaster()
	sqlite.Exec("delete from sqlite_master where type in ('table', 'index', 'trigger');", "deleteAll")
	go sqlite.RegionRegister(connMaster.LocalAddr().String())
	fmt.Println(connMaster.LocalAddr())
	defer connMaster.Close()
	go input(connMaster, StatementChannel)
	go handle(StatementChannel, OutputChannel)
	go output(connMaster, OutputChannel)
	for {
	}

}

func handle(input chan receive, output chan result) {
	for {
		select {
		case rec := <-input:
			fmt.Println("rec", rec)
			var (
				res result
				msg []map[string]interface{}
				err error
			)
			res.Message = "ok"
			res.Error = ""
			switch rec.sqlType {
			case copyStatement:
				res.Message = "copy"
				fmt.Println("复制表:", rec.tableName)
				// 复制表
				res.File, err = getTableLog(rec.tableName)
			case queryStatement:
				fmt.Println("查询语句", rec.sqlStatement)
				msg, err = sqlite.Query(rec.sqlStatement)
			case nonQueryStatement:
				fmt.Println("执行语句", rec.sqlStatement)
				//msg, err = sqlite.Exec(rec.sqlStatement)
				msg, err = sqlite.Exec(rec.sqlStatement, rec.tableName)
			}
			//fmt.Println("sqlExec", msg, err)
			if err != nil {
				fmt.Println(">>> 出现错误!", err)
				res.Error = err.Error()
				res.Message = "NOT OK"
			}
			res.Data = msg
			res.TableList = getTableList()
			res.ClientIP = rec.ipAddress
			fmt.Println(">>> res.ClientIP", res.ClientIP)
			fmt.Println(">>> 得到结果: ", res)
			// 返回结果给 master
			output <- res
		}
	}
}

func input(connMaster net.Conn, input chan receive) {
	//通信得到结果
	for {
		msg := make([]byte, 255)
		msgRead, err := connMaster.Read(msg)
		fmt.Println(">>> msgRead: ", msgRead)
		data := make([]byte, msgRead)
		copy(data, msg)
		if msgRead == 0 || err != nil {
			panic(err)
		} else {
			request := regionRequest{
				TableName: "",
				IpAddress: "",
				Kind:      "",
				Sql:       "",
				File:      make([]string, 0),
			}
			json.Unmarshal(data, &request)
			fmt.Println(">>> request.IpAddress", request.IpAddress)
			fmt.Println(">>> 收到请求: ", request.IpAddress, request.Kind, request.Sql)
			//fmt.Println("killall")
			if request.Kind == "copy" {
				temp := &receive{
					sqlStatement: "",
					sqlType:      copyStatement,
					tableName:    request.TableName,
					ipAddress:    request.IpAddress,
				}
				input <- *temp
			} else if request.Kind == "select" {
				temp := &receive{
					sqlStatement: request.Sql,
					sqlType:      queryStatement,
					tableName:    request.TableName,
					ipAddress:    request.IpAddress,
				}
				input <- *temp
			} else {
				temp := &receive{
					sqlStatement: request.Sql,
					sqlType:      nonQueryStatement,
					tableName:    request.TableName,
					ipAddress:    request.IpAddress,
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
				log.Println(err)
			}
		}
	}
}

// 得到当前db的tablelist
func getTableList() (res []string) {
	sql := "select * from sqlite_master where type = \"table\""
	rawData, _ := sqlite.Query(sql)
	for _, val := range rawData {
		res = append(res, val["name"].(string))
	}
	fmt.Println("当前Table", res)
	return res
}

func getTableLog(tableName string) (res []string, err error) {
	file, err := os.OpenFile("./backup/"+tableName+".txt", os.O_RDONLY, 0644)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer file.Close()
	buf := bufio.NewScanner(file)
	for {
		if !buf.Scan() {
			break
		}
		line := buf.Text()
		res = append(res, line)
	}
	return res, nil
}
