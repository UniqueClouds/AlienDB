package master

import (
	"fmt"
	"sync"
)

var lock sync.Mutex

type copyTableInfo struct {
	tableName   string
	aliveRegion string
}

type tableInfo struct {
	tableName string
	region_1  string
	region_2  string
}

type tableList []*tableInfo

func (tq tableList) Len() int { return len(tq) }

func (tq tableList) Find(tableName string) int {
	index := 0
	for ; index < tq.Len(); index++ {
		if tq[index].tableName == tableName {
			return index
		}
	}
	return -1
}

func (tq tableList) getRegionNumber(tableName string) int {
	index := tq.Find(tableName)
	if index >= 0 {
		if tq[index].region_1 == "" && tq[index].region_2 == "" {
			return 0
		} else if tq[index].region_1 == "" && tq[index].region_2 != "" {
			return 1
		} else if tq[index].region_1 != "" && tq[index].region_2 == "" {
			return 1
		} else {
			return 2
		}
	}
	return -1
}

func (tq tableList) getRegionIp(tableName string) (string, string) {
	region_1, region_2 := "", ""
	if index := tq.Find(tableName); index >= 0 {
		region_1, region_2 = tq[index].region_1, tq[index].region_2
	}
	return region_1, region_2
}

func (tq tableList) updateRegionIp(tableName, ip string) {
	lock.Lock()
	defer lock.Unlock()
	if index := tq.Find(tableName); index >= 0 {
		if tq[index].region_1 == "" && tq[index].region_2 == "" {
			tq[index].region_1 = ip
		} else if tq[index].region_1 == "" && tq[index].region_2 != ip {
			tq[index].region_1 = ip
		} else if tq[index].region_2 == "" && tq[index].region_1 != ip {
			tq[index].region_2 = ip
		}
	} else {
		newTable := tableInfo{
			tableName: tableName,
			region_1:  ip,
			region_2:  "",
		}
		tableQueue = append(tableQueue, &newTable)
	}
}

func (tq tableList) downRegionIp(ip string) []copyTableInfo {
	var tableNeedCopy []copyTableInfo
	for i := 0; i < tq.Len(); i++ {
		if tq[i].region_1 == ip {
			tq[i].region_1 = ""
			info := copyTableInfo{
				tableName:   tq[i].tableName,
				aliveRegion: tq[i].region_2,
			}
			tableNeedCopy = append(tableNeedCopy, info)
		} else if tq[i].region_2 == ip {
			tq[i].region_2 = ""
			info := copyTableInfo{
				tableName:   tq[i].tableName,
				aliveRegion: tq[i].region_1,
			}
			tableNeedCopy = append(tableNeedCopy, info)
		}
	}
	return tableNeedCopy
}

func (tq tableList) downTableIp(tableName, ip string) {
	index := tq.Find(tableName)
	if index == -1 {
		fmt.Printf("> Master: There is no table(%s).\n", tableName)
	} else if tq[index].region_1 == ip {
		tq[index].region_1 = ""
	} else if tq[index].region_2 == ip {
		tq[index].region_2 = ""
	}
}

func (tq tableList) getSameIP(table_1, table_2 string) string {
	ip_1, ip_2 := tq.getRegionIp(table_1)
	ip_3, ip_4 := tq.getRegionIp(table_2)
	if ip_1 == ip_3 || ip_1 == ip_4 {
		return ip_1
	}
	return ip_2
}

func test_2() {
	tableQueue.updateRegionIp("hh", "123")
	tableQueue.updateRegionIp("hh", "123")
	tableQueue.updateRegionIp("hh", "456")

	tableQueue.updateRegionIp("haha", "789")

	table := tableQueue.downRegionIp("123")

	fmt.Println(tableQueue[0].region_1, tableQueue[0].region_2)
	fmt.Println(table[0].aliveRegion, table[0].tableName)
}
