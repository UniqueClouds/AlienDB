package master

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

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

func ListenClient() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":2379")
	if err != nil {
		fmt.Println("client net.tcpAddr error : ", err)
		return
	}
	listen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("client net.Listen error : ", err)
		return
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("client listen.Accept error : ", err)
			continue
		}
		go handleClientRequest(conn)
	}
}

// 接受从节点tcp连接
// func startRegionService() {
// 	// 监听从节点tcp连接
// 	listen, err := net.Listen("tcp", "127.0.0.1:2380")
// 	if err != nil {
// 		fmt.Println("region net.Listen  error : ", err)
// 		return
// 	}

// 	for {
// 		// 建立从节点连接
// 		conn, err := listen.Accept()
// 		if err != nil {
// 			fmt.Println("region listen.Accept error : ", err)
// 			continue
// 		}
// 		// 专门开一个goroutine去处理连接
// 		go startRegion(conn)
// 	}
// }

// 处理客户端连接
func handleClientRequest(connClient net.Conn) {
	fmt.Println("Client " + connClient.RemoteAddr().String() + " Connected!!!")
	defer connClient.Close()
	for {
		msg := make([]byte, 255)
		msgRead, err := connClient.Read(msg)
		data := make([]byte, msgRead)
		copy(data, msg)
		if msgRead == 0 || err != nil {
			panic(err)
		} else {
			request := make(map[string]string)
			json.Unmarshal(data, &request)
			if request["error"] == "" {
				fmt.Println(request["name"], request["kind"], request["sql"], request["join"])
				if request["join"] == "true" {
					continue;
				} else if request["kind"] == "create" {
					handleCreate(request["name"], request["sql"])
				} else {
					handleOther(request["name"], request["sql"])
				}
				msg := map[string]string {"error":""}
				b_msg, _ := json.Marshal(msg)
				connClient.Write(b_msg)
			} else {
				msg := map[string]string {"error":"Illegal request."}
				b_msg, _ := json.Marshal(msg)
				connClient.Write(b_msg)
			}
		}
	}
}
// 处理从节点连接
// func startRegion(conn net.Conn) {
// 	defer conn.Close()
// 	// 新建map记录新连接的从节点
// 	// 容错容灾等
// }
