package sqlite

import (
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

func ConnectToMaster() net.Conn {
	var ip, port string
	fmt.Println(">>>请输入目标服务器IP：")
	_, err := fmt.Scan(&ip)
	if err != nil {
		panic(err)
	}
	fmt.Println(">>>请输入目标服务器端口号：")
	_, err = fmt.Scan(&port)
	if err != nil {
		panic(err)
	}
	//var connMaster net.Conn
	//ip = "localhost"
	//port = "2223"
	connMaster, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(">>> 与 Master 连接成功！")
	}
	return connMaster
}

type Result struct {
	Error string              `json:"error"`
	Data  []map[string]string `json:"data"`
}
