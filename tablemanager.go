package main

type tableInfo struct {
	tableName 			string
	region_1			string
	region_2			string
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

func (tq tableList) getRegionIp(tableName string) (string, string){
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
		newTable := tableInfo {
			tableName: tableName,
			region_1: ip,
			region_2: "",
		}
		tableQueue = append(tableQueue, &newTable)
	}
}