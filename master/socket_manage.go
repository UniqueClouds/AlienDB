package master

import (
	"encoding/json"
	"fmt"
	"net"
)

// 获取本机ip地址，方便客户端及从节点的连接
func GetLocalIP() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("get local ip error : ", err)
		return
	}
	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ipAddr.IP.IsLoopback() {
			continue
		}
		if !ipAddr.IP.IsGlobalUnicast() {
			continue
		}
		fmt.Println("local ip: ", ipAddr.IP.String())
	}
	return
}

// 接受客户端tcp连接
func startClientService() {
	// 监听客户端tcp连接
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":2379")
	if err != nil {
		fmt.Println("client net.tcpAddr error : ", err)
		return
	}

	listen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("client net.Listen  error : ", err)
		return
	}

	for {
		// 建立客户端连接
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("client listen.Accept error : ", err)
			continue
		}
		// 专门开一个goroutine去处理连接
		go startClient(conn)
	}
}

// 接受从节点tcp连接
func startRegionService() {
	// 监听从节点tcp连接
	listen, err = net.Listen("tcp", "127.0.0.1:2380")
	if err != nil {
		fmt.Println("region net.Listen  error : ", err)
		return
	}

	for {
		// 建立从节点连接
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("region listen.Accept error : ", err)
			continue
		}
		// 专门开一个goroutine去处理连接
		go startRegion(conn)
	}
}

// 处理客户端连接
func startClient(connClient net.Conn) {
	defer connClient.Close()
	data := make([]byte, 255)
	msgRead, err := connClient.Read(data)
	if msgRead == 0 || err != nil {
		panic(err)
	} else {
		request := make(map[string]string)
		//字符串转map
		err = json.Unmarshal(data, &request)
		if request["error"] == "" {
			tableName = request["name"]
			actionKind = request["kind"]
			// 后续对表进行处理

			_, err := connClient.Write(msgStr) // 最后返回查询结果
		} else {
			fmt.Println(">>>异常：" + result["error"])
		}
	}
}

// 处理从节点连接
func startRegion(conn net.TCPConn) {
	defer conn.Close()
	// 新建map记录新连接的从节点
	// 容错容灾等
}
