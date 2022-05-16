package ClientManager

import "fmt"

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
