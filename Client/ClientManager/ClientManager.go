package ClientManager

import (
	"fmt"
	"strings"
)

func run() {
	fmt.Print("Distributed-MiniSQL 客户端启动......")
	var input string
	line := ""
	fmt.Scan(&input)
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
