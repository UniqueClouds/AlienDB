package client

import "strings"

// 解释器，简单解析sql语句，查明语句类型和表名
func interpreter(sql string) map[string]string {
	result := make(map[string]string)
	strings.ReplaceAll(sql, "\\s+", " ")
	words := strings.Split(sql, " ")
	result["kind"] = words[0]
	if words[0] == "create" {
		result["name"] = words[2]
	} else if words[0] == "update" || words[0] == "drop" || words[0] == "insert" || words[0] == "delete" {
		name := strings.Trim(words[2], "();")
		result["name"] = name
	} else if words[0] == "select" {
		var i int
		for i = 0; i < len(words); i++ {
			if words[i] == "from" && i != len(words)-1 {
				result["name"] = words[i+1]
			} else if words[i] == "join" && i != len(words)-1 {
				result["name"] = result["name"] + " " + words[i+1]
			}
		}
	}
	result["sql"] = sql
	if _, ok := result["name"]; !ok {
		result["error"] = "true"
	}
	return result
}
