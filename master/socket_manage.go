package master

import (
	"encoding/json"
	"fmt"
	"net"
)

// GetLocalIP 获取本机ip地址，方便客户端及从节点的连接
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("get local ip error : ", err)
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	panic("unable to determine local ip!")
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
func startRegionService(localIP string) {
	// 监听从节点tcp连接
	localIP = localIP + ":8000"
	listen, err = net.Listen("tcp", localIP)
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
