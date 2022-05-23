package master

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

//返回结果，与Region统一
type result struct {
	Error string                   `json:"error"`
	Data  []map[string]interface{} `json:"data"`
}

// GetLocalIP 获取本机ip地址，方便客户端及从节点的连接
func GetLocalIP() string {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if addr.IsLoopback() {
			continue
		}
		if ipv4 := addr.To4(); ipv4 != nil {
			return ipv4.String()
		}
	}
	return ""
}

// ListenRegion 监听Region， 端口号：2223
func ListenRegion() {
	localIP := GetLocalIP()
	fmt.Println(">>> 本机IP", localIP)
	fmt.Println(">>> Master 监听 Region 端口: 2223")
	tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:2223")
	if err != nil {
		fmt.Println(">>> Region net.tcpAddr error : ", err)
		return
	}
	listen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println(">>> Region net.tcpAddr error : ", err)
		return
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(">>> Region net.tcpAddr error : ", err)
			continue
		}
		go connectToRegion(conn)
	}
}

// 处理Region连接
func connectToRegion(connRegion net.Conn) {
	//defer connRegion.Close()
	newRegion := &IpAddressInfo{
		ipAddress: connRegion.RemoteAddr().String(),
		index:     regionList.Len() + 1,
	}
	fmt.Println(">>> 新增Region ...")
	regionList.Push(newRegion)
	ReceiveChan := make(chan result, 500)
	SendChan := make(chan map[string]string, 500)

	request := make(map[string]string)
	// example
	request["kind"] = "select"
	request["sql"] = "SELECT * FROM userinfo"
	SendChan <- request

	fmt.Println(">>> Region (地址" + connRegion.RemoteAddr().String() + ") Connected!!!")
	// 发送消息
	go send(connRegion, SendChan)

	//go handle(ReceiveChan)
	// 接收消息
	go receive(connRegion, ReceiveChan)
}

func send(connRegion net.Conn, SendChan chan map[string]string) {
	sendMsg := <-SendChan
	fmt.Println(">>> 发送给Region 消息: ", sendMsg)
	msgStr, _ := json.Marshal(sendMsg)
	if _, err := connRegion.Write(msgStr); err != nil {
		panic(err)
	}
}

func receive(connRegion net.Conn, ReceiveChan chan result) {
	//通信得到结果
	msg := make([]byte, 255)
	msgRead, err := connRegion.Read(msg)
	data := make([]byte, msgRead)
	copy(data, msg)
	if msgRead == 0 || err != nil {
		panic(err)
	} else {
		ans := result{
			Error: "",
			Data:  nil,
		}
		json.Unmarshal(data, &ans)
		if ans.Error == "" {
			fmt.Println(">>> 收到 Region 消息: ", ans.Data)
		}
	}
}
