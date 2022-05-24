package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type regionRequest struct {
	TableName string
	IpAddress string
	Kind string
	Sql string
}

type clientResult struct {
	Error string                   `json:"error"`
	Data  []map[string]interface{} `json:"data"`
}

type regionResult struct {
	Error string                   `json:"error"`
	Data  []map[string]interface{} `json:"data"`
	TableList []string 			   `json:"tableList"`
	Message string				   `json:"message"`
	ClientIP string 			   `json:"clientIp"`
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

func ListenClient() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":2224")
	if err != nil {
		fmt.Println("client net.tcpAddr error : ", err)
		return
	}
	listen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("client net.Listen error : ", err)
		return
	}
	fmt.Println("> Master: The port listening to client connection is: 2224")
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("client listen.Accept error : ", err)
			continue
		}
		go sessionWithClient(conn)
	}
}

func ListenRegion() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":2223")
	if err != nil {
		fmt.Println(">>> Region net.tcpAddr error : ", err)
		return
	}
	listen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println(">>> Region net.tcpAddr error : ", err)
		return
	}
	fmt.Println("> Master: The port listening to region connection is: 2223")
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(">>> Region net.tcpAddr error : ", err)
			continue
		}
		go sessionWithRegion(conn)
	}
}

func sessionWithClient(connClient net.Conn) {
	fmt.Println("> Master: Client " + connClient.RemoteAddr().String() + " Connected.")
	newClient := &clientInfo{
		ipAddress: connClient.RemoteAddr().String(),
		resultQueue:	make(chan clientResult, 20),
	}
	clientQueue = append(clientQueue, newClient)
	defer connClient.Close()
	go handleClientRequest(connClient)
	go sendClientResult(connClient)
	select {}
}

func handleClientRequest(connClient net.Conn) {
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
				if request["join"] == "true" {
					
				} else if request["kind"] == "create" {
					handleCreate(connClient.RemoteAddr().String(), request["name"], request["sql"])
				} else {
					handleOther(connClient.RemoteAddr().String(), request["name"], request["kind"], request["sql"])
				}
			} else {

			}
		}
	}
}

func sendClientResult(connClient net.Conn) {
	for {
		if id := clientQueue.Find(connClient.RemoteAddr().String()); id >= 0 {
			select {
				case sendRlt := <-clientQueue[id].resultQueue :
					fmt.Printf("> Master: Send to client(%s) [%s].\n", connClient.RemoteAddr().String(), sendRlt.Data)
					msgStr, _ := json.Marshal(sendRlt)
					if _, err := connClient.Write(msgStr); err != nil {
						panic(err)
				}
			}
		}
	}
}

func sessionWithRegion(connRegion net.Conn) {
	fmt.Println("> Master: Region " + connRegion.RemoteAddr().String() + " Connected.")
	newRegion := &IpAddressInfo{
		ipAddress: connRegion.RemoteAddr().String(),
		requestQueue:	make(chan regionRequest, 20),
		tableNumber: 0,
	}
	regionQueue.Push(newRegion)
	fmt.Printf("> Master: There are now %d region connections.\n", regionQueue.Len())
	go sendRegionRequest(connRegion)
	go handleRegionReceive(connRegion)
	select {}
}

func sendRegionRequest(connRegion net.Conn) {
	for {
		if id := regionQueue.find(connRegion.RemoteAddr().String()); id >= 0 {
			select {
				case sendMsg := <-regionQueue[id].requestQueue :
					msgStr, _ := json.Marshal(sendMsg)
					fmt.Printf("> Master: Send to region(%s) [%s].\n", connRegion.RemoteAddr().String(), sendMsg.Sql)
					if _, err := connRegion.Write(msgStr); err != nil {
						panic(err)
				}
			}
		}
	}
}

func handleRegionReceive(connRegion net.Conn) {
	for {
		msg := make([]byte, 255)
		msgRead, err := connRegion.Read(msg)
		data := make([]byte, msgRead)
		copy(data, msg)
		if msgRead == 0 || err != nil {
			panic(err)
		} else {
			rec := regionResult{
				Error: "",
				Data:  nil,
				Message: "",
				ClientIP: "",
			}
			json.Unmarshal(data, &rec)
			regionQueue.update(connRegion.RemoteAddr().String(), len(rec.TableList))
			for _, tn := range rec.TableList {
				tableQueue.updateRegionIp(tn, connRegion.RemoteAddr().String())
			}
			if rec.Error == "" {
				fmt.Printf("> Master: Receive data from region(%s): %s\n", connRegion.RemoteAddr().String(), rec.Message)
				handleResult(rec)
			} else {
				fmt.Println("Error: ", rec.Error)
			}
		}
	}
}
