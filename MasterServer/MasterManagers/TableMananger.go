package MasterManagers

import (
	"container/list"
	"net"
)

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
