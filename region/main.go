package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"my/AlienDB/region/sqlite"
	"net"
)

var QuitChan chan string

type receive struct {
	sqlStatement string
	sqlType      int
	ipAddress    string
}

type regionRequest struct {
	IpAddress string
	Kind      string
	Sql       string
}
type result struct {
	Error     string                   `json:"error"`
	Data      []map[string]interface{} `json:"data"`
	TableList []string                 `json:"tableList"`
	Message   string                   `json:"message"`
	ClientIP  string                   `json:"clientIP"`
}

const (
	queryStatement    = 1
	nonQueryStatement = 2
	quitStatement     = 3
)

func main() {
	defer sqlite.Close()
	fmt.Println(">>> Region 启动中...")
	//func() {
	//
	//	//listener, err := net.ListenTCP("tcp", "localhost:2379")
	//	var endpoints = []string{"localhost:2223"}
	//	// 暂定名称
	//	ser, err := sqlite.NewServiceRegister(endpoints, "/db/region_01", "localhost:8000", 5)
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//	// 监听续租相应 chan
	//	go ser.ListenLeaseRespChan()
	//
	//	// 监听系统信号，等待 ctrl + c 系统信号通知服务关闭
	//	//c := make(chan os.Signal, 1)
	//	select {}
	//}()
	//return
	//var endpoints = []string{"localhost:2222"}
	StatementChannel := make(chan receive, 500)
	OutputChannel := make(chan result, 500)
	QuitChan = make(chan string)
	//go func() {
	//	ser, err := sqlite.NewServiceRegister(endpoints, "/db/region_01", "localhost:8000", 6, 5)
	//	if err != nil {
	//		fmt.Println(err)
	//		log.Fatalln(err)
	//	}
	//	go ser.ListenLeaseRespChan()
	//}()
	connMaster := sqlite.ConnectToMaster()
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

			switch rec.sqlType {
			case queryStatement:
				fmt.Println("查询语句", rec.sqlStatement)
				msg, err = sqlite.Query(rec.sqlStatement)

			case nonQueryStatement:
				fmt.Println("执行语句", rec.sqlStatement)
				msg, err = sqlite.Exec(rec.sqlStatement)
				msg, err = sqlite.Exec(rec.sqlStatement, "tablename")
			}
			//fmt.Println("sqlExec", msg, err)
			res.Error = ""
			res.Message = "OK"
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
			if err != nil {
				QuitChan <- "err"
			} else {
				output <- res
			}
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
				IpAddress: "",
				Kind:      "",
				Sql:       "",
			}
			json.Unmarshal(data, &request)
			fmt.Println(">>> request.IpAddress", request.IpAddress)
			fmt.Println(">>> 收到请求: ", request.IpAddress, request.Kind, request.Sql)
			fmt.Println("killall")
			if request.Kind == "select" {
				temp := &receive{
					sqlStatement: request.Sql,
					sqlType:      queryStatement,
					ipAddress:    request.IpAddress,
				}
				input <- *temp
			} else {
				temp := &receive{
					sqlStatement: request.Sql,
					sqlType:      nonQueryStatement,
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
				panic(err)
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
