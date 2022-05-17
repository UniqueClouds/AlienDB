package master

import "net"

/* 一个用于记录各种信息的表
* 需要如下结构
*
* list 记录所有连接过的从节点ip
* map 映射当前活跃的从节点ip与每个ip对应的tablelist
* 功能:
1. 表的增加删除 addTable, deleteTable
2. 根据表获得ip地址
3. 负载均衡
*/
type tableManager struct {
}

type TableManager struct {
	tableInfo       map[string]string
	serverList      list.List
	aliveServer     map[string]list.List
	socketThreadMap map[string]net.Conn
}

func initTableManager() {
	tableManager := TableManager{tableInfo: make(map[string]string), aliveServer: make(map[string]list.List), socketThreadMap: make(map[string]net.Conn)}
}

func (tableManager TableManager) addTable(table string, inetAddress string) {
	tableManager.tableInfo[table] = inetAddress
	if _, ok := tableManager.aliveServer[inetAddress]; ok {
		l := tableManager.aliveServer[inetAddress]
		l.PushBack(table)
	} else {
		temp := list.List{}
		temp.Init()
		temp.PushBack(table)
		tableManager.aliveServer[inetAddress] = temp
	}
}

func (tableManager TableManager) deleteTable(table string, inetAddress string) {
	delete(tableManager.tableInfo, table)
	l := tableManager.aliveServer[inetAddress]
	l.Remove(table)
}
