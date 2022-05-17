package client

import (
	"fmt"
	"net"
	"strings"
)

func RunClient() {
	fmt.Print("TonyDB initiating......")
	var c Cache
	c.initCache()
	line := ""
	for {
		var sql string
		fmt.Println("新消息>>>请输入你想执行的SQL语句：")
		for len(line) == 0 || line[len(line)-1] != ';' {
			fmt.Scanln(&line)
			if len(line) == 0 {
				continue
			}
			sql = sql + line
			sql += " "
		}
		line = ""
		fmt.Println(sql)
		if strings.Contains(sql, "quit;") {
			// 退出
			// 关闭socket 服务器
			break

		}
		command := sql
		sql = ""
		target := interpreter(command)
		if _, ok := target["error"]; ok {
			fmt.Println("新消息>>>输入有误，请重试!")
		}
		table := target["name"]
		var cache string = ""
		fmt.Println("新消息>>>需要处理的表名是: " + table)
		/*******
		待 完 善
		*/
		if strings.Compare(target["kind"], "create") == 0 {
			// 处理建表
		} else {
			if strings.Compare(target["Cache"], "true") == 0 {
				cache = c.getCache(table)
				if cache == "" {
					fmt.Println("新消息>>>客户端缓存中不存在该表!")
				} else {
					fmt.Println("新消息>>>客户端缓存中存在该表！其对应的服务器是：" + cache)
				}
				if cache == "" {
					// 到masterSocket里面查询
					// 缺
				} else {
					//找到了就直接连接 regionsocketmanager
					// 缺
				}
			}
		}
	}
}

func connectToMaster() {
	connMaster, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		panic(err)
	}
	defer connMaster.Close()
	for {
		msg := ""
		fmt.Println("发送给Master:", msg)
		connMaster.Write([]byte(msg))

		data := make([]byte, 255)
		msg_read, err := connMaster.Read(data)
		if msg_read == 0 || err != nil {
			panic(err)
		}
		msg_read_str := string(data[0:msg_read])
		if msg_read_str == "close" {
			connMaster.Write([]byte("close"))
			break
		}
	}

}

func connectToRegion(sql string, regionIPAddr string) {
	fmt.Println("新消息>>>与从节点" + regionIPAddr + "建立连接")
	connRegion, err := net.Dial("tcp", regionIPAddr)
	if err != nil {
		panic(err)
	}
	defer connRegion.Close()
	for {
		msg := sql
		fmt.Println("发送给Region:", msg)
		connRegion.Write([]byte(msg))
		data := make([]byte, 255)
		msg_read, err := connRegion.Read(data)
		if msg_read == 0 || err != nil {
			panic(err)
		}

	}
}
