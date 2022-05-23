package master

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

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

func ListenRegion() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "8000")
	if err != nil {
		fmt.Println("Region net.tcpAddr error : ", err)
		return
	}
	listen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("Region net.tcpAddr error : ", err)
		return
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Region net.tcpAddr error : ", err)
			continue
		}
		go connectToRegion(conn)
	}
}

type result struct {
	Error string                   `json:"error"`
	Data  []map[string]interface{} `json:"data"`
}

// 处理Region连接
func connectToRegion(connRegion net.Conn) {
	ReceiveChan := make(chan result, 500)
	SendChan := make(chan result, 500)
	//QuitChan := make(chan string)
	fmt.Println(">>> Region " + connRegion.RemoteAddr().String() + " Connected!!!")
	defer connRegion.Close()
	for {
		go send(connRegion, SendChan)
		//go handle(ReceiveChan)
		go receive(connRegion, ReceiveChan)
	}
}

func send(connRegion net.Conn, SendChan chan result) {
	fmt.Println(">>> 发送给Region 消息 ...")
	sendMsg := <-SendChan
	msgStr, _ := json.Marshal(sendMsg)
	fmt.Println(msgStr)
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
		request := make(map[string]string)
		if request["error"] == "" {
			fmt.Println(">>> 收到Region消息:")
			//
			fmt.Println(request)
		}
	}
}
