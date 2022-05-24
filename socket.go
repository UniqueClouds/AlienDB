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
	File []string
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
	File []string				   `json:"file"`
}

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
		resultQueue: make(chan clientResult, 20),
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
			// panic(err)
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
			var msg clientResult
			for i := 0; i < 2; i++ {
				select {
					case sendRlt := <-clientQueue[id].resultQueue :
						if i == 0 {
							msg = sendRlt
						} else if sendRlt.Error == "" {
							msg = sendRlt
						}
				}
			}
			fmt.Printf("> Master: Send to client(%s) [%s].\n", connClient.RemoteAddr().String(), msg.Data)
			msgStr, _ := json.Marshal(msg)
			if _, err := connClient.Write(msgStr); err != nil {
				panic(err)
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

	go copyRequest(connRegion)

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
			// panic(err)
		} else {
			rec := regionResult{
				Error: "",
				Data:  nil,
				Message: "",
				ClientIP: "",
				File: nil,
			}
			json.Unmarshal(data, &rec)
			regionQueue.update(connRegion.RemoteAddr().String(), len(rec.TableList))
			for _, tn := range rec.TableList {
				tableQueue.updateRegionIp(tn, connRegion.RemoteAddr().String())
			}
			if rec.Error == "" {
				fmt.Printf("> Master: Receive data from region(%s): %s\n", connRegion.RemoteAddr().String(), rec.Message)
				if rec.Message == "copy" {
					forwardCopy(rec, connRegion)
				} else {
					handleResult(rec)
				}
			} else {
				handleError(rec)
			}
		}
	}
}

func copyRequest(conn net.Conn) {
	id := regionQueue.find(conn.RemoteAddr().String())
	for {
			select {
				case tableName := <-regionQueue[id].copyRequestQueue  :
					request := regionRequest{
						TableName: tableName,
						IpAddress: "",
						Kind: "copy",
						Sql: "",
						File: nil,
					}
					fmt.Printf("> Master: Send to client(%s) [copy %s].\n", conn.RemoteAddr().String(), request.TableName)
					msgStr, _ := json.Marshal(request)
					if _, err := conn.Write(msgStr); err != nil {
						panic(err)
				}
			}
		}
}

func forwardCopy(rec regionResult, conn net.Conn) {
	desRegion := regionQueue.getCopyRegion(conn.RemoteAddr().String())
	request := regionRequest {
		TableName: "",
		IpAddress: "",
		Kind: "copy",
		Sql: "",
		File: rec.File,
	}
	regionQueue[regionQueue.find(desRegion)].requestQueue <- request
}