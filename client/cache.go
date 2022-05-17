package client

import "fmt"

// 客户端缓存管理

type Cache struct {
	cache map[string]string
}

// 缓存表
//var Cache map[string]string

// 初始化分配内存
func (c Cache) initCache() {
	c.cache = make(map[string]string)
}

// 查询缓存
func (c Cache) getCache(table string) string {
	if _, ok := c.cache[table]; ok {
		return c.cache[table]
	}
	return ""
}

// 存储已知的表和对应的服务器
func (c Cache) setCache(table, server string) {
	c.cache[table] = server
	fmt.Print("存入缓存：表明" + table + "端口号" + table)
}
