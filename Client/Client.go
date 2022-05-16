package Client

import (
	"fmt"
	"net"
	"strings"
)

// 客户端缓存管理

// 缓存表
var cache map[string]string

// 初始化分配内存
func init() {
	cache = make(map[string]string)
}

// 查询缓存
func getCache(table string) string {
	if _, ok := cache[table]; ok {
		return cache[table]
	}
	return ""
}

// 存储已知的表和对应的服务器
func setCache(table, server string) {
	cache[table] = server
	fmt.Print("存入缓存：表明" + table + "端口号" + table)
}
func RunClient() {
	fmt.Print("TonyDB initiating......")
	line := ""
	for true {
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
			if strings.Compare(target["cache"], "true") == 0 {
				cache = getCache(table)
				if cache == "" {
					fmt.Println("新消息>>>客户端缓存中不存在该表!")
				} else {
					fmt.Println("新消息>>>客户端缓存中存在该表！其对应的服务器是：" + cache)
				}
				if cache == "" {
					// 到masterSocket里面查询
				} else {
					//找到了就直接连接 regionsocketmanager
				}
			}
		}
	}
}

func interpreter(sql string) map[string]string {
	var result map[string]string = make(map[string]string)
	result["cache"] = "true"
	strings.ReplaceAll(sql, "\\s+", " ")
	words := strings.Split(sql, " ")
	result["kind"] = words[0]
	if strings.Compare(words[0], "create") == 0 {
		result["cache"] = "false"
		result["name"] = words[2]
	} else if strings.Compare(words[0], "drop") == 0 || strings.Compare(words[0], "insert") == 0 || strings.Compare(words[0], "delete") == 0 {
		name := strings.Trim(words[2], "();")
		result["name"] = name
	} else if strings.Compare(words[0], "select") == 0 {
		for i := 0; i < len(words); i++ {
			if strings.Compare(words[i], "from") == 0 && i != len(words)-1 {
				result["name"] = words[i+1]
				break
			}
		}
	}
	if _, ok := result["name"]; !ok {
		result["error"] = "true"
	}
	return result
}

func ConnectToMaster() {
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
func ConnectToRegion(sql string, regionIPAddr string) {
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
