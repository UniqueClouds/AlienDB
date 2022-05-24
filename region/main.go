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
	TableName string                   `json:"tableName"`
}

const (
	queryStatement    = 1
	nonQueryStatement = 2
	quitStatement     = 3
	copyStatement     = 4
	dropStatement     = 5
	newStatement      = 6
)

func main() {
	defer sqlite.Close()
	StatementChannel := make(chan receive, 500)
	OutputChannel := make(chan result, 500)
	QuitChan = make(chan string)
	fmt.Println(">Region: 启动中...")
	connMaster := sqlite.ConnectToMaster()
	//sqlite.Exec("delete from sqlite_master where type in ('table', 'index', 'trigger');", "deleteAll")
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
			fmt.Println("> Region: rec", rec)
			var (
				res result
				msg []map[string]interface{}
				err error
			)
			res.Message = "ok"
			res.Error = ""
			switch rec.sqlType {
			case newStatement:
				fmt.Println("> Region: 还原表: ", rec.tableName)
				fmt.Println("> Region: file []string: ", rec.file)
				for _, query := range rec.file {
					fmt.Println("> Region: 执行语句: ", query)
					msg, err = sqlite.Exec(query, rec.tableName)
					if err != nil {
						fmt.Println("err", err)
					}
				}
				res.Message = "copy ok"
			case copyStatement:
				res.Message = "copy"
				fmt.Println("> Region: 复制表: ", rec.tableName)
				// 复制表
				res.File, err = getTableLog(rec.tableName)
			case queryStatement:
				fmt.Println("> Region: 查询语句", rec.sqlStatement)
				msg, err = sqlite.Query(rec.sqlStatement)
			case dropStatement:
				fmt.Println("> Region: Drop语句", rec.sqlStatement)
				msg, err = sqlite.Exec(rec.sqlStatement, rec.tableName)
				res.Message = "drop ok"
			case nonQueryStatement:
				fmt.Println("> Region: 执行语句", rec.sqlStatement)
				//msg, err = sqlite.Exec(rec.sqlStatement)
				msg, err = sqlite.Exec(rec.sqlStatement, rec.tableName)
			}
			//fmt.Println("sqlExec", msg, err)
			if err != nil {
				fmt.Println("> Region:  出现错误: ", err)
				res.Error = err.Error()
				res.Message = "NOT OK"
			}
			res.Data = msg
			res.TableList = getTableList()
			res.ClientIP = rec.ipAddress
			res.TableName = rec.tableName
			fmt.Println("> Region:  res.ClientIP", res.ClientIP)
			fmt.Println("> Region:  得到结果: ", res)
			// 返回结果给 master
			output <- res
		}
	}
}

func input(connMaster net.Conn, input chan receive) {
	//通信得到结果
	for {
		msg := make([]byte, 1024*10)
		msgRead, err := connMaster.Read(msg)
		fmt.Println("> Region:  msgRead: ", msgRead)
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
				File:      nil,
			}
			json.Unmarshal(data, &request)
			fmt.Println("> Region: 收到请求: ", "ipaddress", request.IpAddress, "kind", request.Kind, "sql", request.Sql, "tableName", request.TableName, "file", request.File)

			if request.Kind == "new" {
				temp := &receive{
					sqlStatement: "",
					sqlType:      newStatement,
					tableName:    request.TableName,
					ipAddress:    request.IpAddress,
					file:         request.File,
				}
				input <- *temp
			} else if request.Kind == "copy" {
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
			} else if request.Kind == "drop" {
				temp := &receive{
					sqlStatement: request.Sql,
					sqlType:      dropStatement,
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

			fmt.Println("> Region: 返回给 Master: ", outPutMsg)
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
	fmt.Println("> Region: 当前Table", res)
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
