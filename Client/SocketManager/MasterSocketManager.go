package SocketManager

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var sig = make(chan os.Signal)

func MasterSocketManager() {
	// 解析服务端地址
	RemoteAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8000")
	if err != nil {
		panic(err)
	}
	// 解析本地连接地址
	LocalAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8001")
	if err != nil {
		panic(err)
	}
	// 连接服务端
	conn, err := net.DialTCP("tcp", LocalAddr, RemoteAddr)
	if err != nil {
		panic(err)
	}
	HandleConnectionForMaster(conn)
}

func HandleConnectionForMaster(conn net.Conn) {
	// 监控系统信号
	go signalMonitor(conn)
	// 初始化一个缓存区
	Stdin := bufio.NewReader(os.Stdin)
	for {
		getResponse(conn)
		fmt.Print("[ random_w ]# ")
		input, err := Stdin.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
	}
}
