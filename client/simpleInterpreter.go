package client

import "strings"

// 解释器，简单解析sql语句，查明语句类型和表名
func interpreter(sql string) map[string]string {
	var result map[string]string = make(map[string]string)
	result["Cache"] = "true"
	strings.ReplaceAll(sql, "\\s+", " ")
	words := strings.Split(sql, " ")
	result["kind"] = words[0]
	if strings.Compare(words[0], "create") == 0 {
		result["Cache"] = "false"
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
