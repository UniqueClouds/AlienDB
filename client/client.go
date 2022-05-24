package client

import (
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strings"
)

func RunClient() {
	fmt.Println("AlienDB初始化中......")
	//别管了
	//var c Cache
	//c.initCache()

	//与主节点建立连接
	connMaster := connectToMaster()
	for {
		var sql, line string
		fmt.Println(">>>请输入你想执行的SQL语句：")
		//逐词读取SQL语句
		for len(line) == 0 || line[len(line)-1] != ';' {
			_, err := fmt.Scan(&line)
			if err != nil {
				panic(err)
			}
			sql += line
			if len(line) != 0 {
				sql += " "
			}
		}
		if strings.Contains(sql, "quit;") {
			// 关闭socket 服务器
			msgStr, _ := json.Marshal(map[string]string{"kind": "quit"})
			_, err := connMaster.Write(msgStr)
			if err != nil {
				panic(err)
			}
			err = connMaster.Close()
			if err != nil {
				panic(err)
			}
			break
		}
		//解析操作类型与表名
		target := Interpreter(sql)
		if _, ok := target["error"]; ok {
			fmt.Println(">>>输入有误，请重试！")
		}
		fmt.Println(">>>需要处理的表名是：" + target["name"])
		//直接向主节点发送请求
		sendToMaster(connMaster, &target)
		//别管了
		//if target["kind"] == "create" {
		//	sendToMaster(connMaster, &target, &c)
		//} else {
		//	cache := c.getCache(target["name"])
		//	if cache == nil {
		//		fmt.Println(">>>客户端缓存中不存在该表！")
		//		sendToMaster(connMaster, &target, &c)
		//		cache = c.getCache(target["name"])
		//	} else {
		//		fmt.Println(">>>客户端缓存中存在该表！")
		//		//fmt.Println(">>>客户端缓存中存在该表！其对应的服务器是：" + cache)
		//	}
		//	go connectToRegion(sql, cache)
		//}
	}
}

// 返回结果，与region统一
type Result struct {
	Error string                   `json:"error"`
	Data  []map[string]interface{} `json:"data"`
}

func sendToMaster(connMaster net.Conn, target *map[string]string) {
	//map转化为字符串
	msgStr, _ := json.Marshal(*target)
	_, err := connMaster.Write(msgStr)
	if err != nil {
		panic(err)
	}
	fmt.Println(">>>消息发送成功！")
	rawData := make([]byte, 255)
	//rawDataCopy := make([]byte, 255)
	msgRead, err := connMaster.Read(rawData)
	//msgReadCopy, errCopy := connMaster.Read(rawDataCopy)
	if msgRead == 0 || err != nil {
		//if !(msgReadCopy == 0 || errCopy != nil) {
		//	rawData = rawDataCopy
		//	msgRead = msgReadCopy
		//} else {
		panic(err)
		//}
	} else {
		var result Result
		data := make([]byte, msgRead)
		copy(data, rawData)
		fmt.Println(">>>收到回复：" + string(data))
		//字符串转map
		err = json.Unmarshal(data, &result)
		if result.Error == "" {
			fmt.Println(">>>操作成功！")
			if (*target)["kind"] == "select" {
				//打印查询结果
				fmt.Println(">>>查询结果如下：")
				table := result.Data
				var col []string
				length := make(map[string]int)
				for i, row := range table {
					if i == 0 {
						for k := range table[0] {
							//fmt.Printf("%10s", k)
							col = append(col, k)
							length[k] = len(k)
						}
						sort.Strings(col)
						//fmt.Println("")
					}
					for _, name := range col {
						strRowName := fmt.Sprint(row[name])
						if len(strRowName) > length[name] {
							length[name] = len(strRowName)
						}
						//fmt.Printf("%10v", row[name])
					}
				}
				for i, name := range col {
					if i == 0 {
						fmt.Print("|")
					}
					for k := 0; k < length[name]-len(name)+1; k++ {
						fmt.Print(" ")
					}
					fmt.Printf("%2s |", name)
				}
				fmt.Println("")
				for _, row := range table {
					for i, name := range col {
						if i == 0 {
							fmt.Print("|")
						}
						for j := 0; j < length[name]-len(fmt.Sprint(row[name]))+1; j++ {
							fmt.Print(" ")
						}
						fmt.Printf("%v |", row[name])
					}
					fmt.Println("")
				}
				//fmt.Println(result["data"])
			}
			//别管了
			//if msg["kind"] == "create" {
			//	fmt.Println(">>>操作成功！")
			//} else {
			//	//cache.setCache(msg["name"], result["server"])
			//	fmt.Println(">>>缓存已更新！")
			//}
		} else {
			fmt.Println(">>>异常：" + result.Error)
		}
	}
}

func connectToMaster() net.Conn {
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
	connMaster, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(">>>连接成功！")
	}
	return connMaster
}

//别管了
//func connectToRegion(sql string, regionIPAddr []string) {
//	fmt.Println(">>>与从节点" + regionIPAddr + "建立连接")
//	connRegion, err := net.Dial("tcp", regionIPAddr)
//	if err != nil {
//		panic(err)
//	}
//	defer func(connRegion net.Conn) {
//		err := connRegion.Close()
//		if err != nil {
//			panic(err)
//		}
//	}(connRegion)
//	for {
//		msg := sql
//		fmt.Println(">>>发送给Region：", msg)
//		_, err := connRegion.Write([]byte(msg))
//		if err != nil {
//			panic(err)
//		}
//		data := make([]byte, 255)
//		msgRead, err := connRegion.Read(data)
//		if msgRead == 0 || err != nil {
//			panic(err)
//		}
//	}
//}
