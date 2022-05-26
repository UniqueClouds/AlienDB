package master

var tableQueue tableList
var clientQueue clientList
var regionQueue PriorityQueue

type middleResult struct {
	result clientResult
	times  int
}

func handleCreate(ipAddress, tableName, sqlInstruction string) {
	desRegion_1, desRegion_2 := regionQueue.getNextTwo()
	request := regionRequest{
		TableName: tableName,
		IpAddress: ipAddress,
		Kind:      "create",
		Sql:       sqlInstruction,
		File:      nil,
	}
	regionQueue[regionQueue.find(desRegion_1)].requestQueue <- request
	regionQueue[regionQueue.find(desRegion_2)].requestQueue <- request
}

func handleOther(ipAddress, tableName, kind, sqlInstruction string) {
	desRegion_1, desRegion_2 := tableQueue.getRegionIp(tableName)
	request := regionRequest{
		TableName: tableName,
		IpAddress: ipAddress,
		Kind:      kind,
		Sql:       sqlInstruction,
		File:      nil,
	}
	if desRegion_1 != "" && desRegion_2 != "" {
		regionQueue[regionQueue.find(desRegion_1)].requestQueue <- request
		regionQueue[regionQueue.find(desRegion_2)].requestQueue <- request
	} else if desRegion_1 != "" && desRegion_2 == "" {
		regionQueue[regionQueue.find(desRegion_1)].requestQueue <- request
	} else if desRegion_2 != "" && desRegion_1 == "" {
		regionQueue[regionQueue.find(desRegion_2)].requestQueue <- request
	}
}

func handleJoin(ipAddress, ip, sqlInstruction string) {
	request := regionRequest{
		TableName: "",
		IpAddress: ipAddress,
		Kind:      "select",
		Sql:       sqlInstruction,
		File:      nil,
	}
	regionQueue[regionQueue.find(ip)].requestQueue <- request
}

func handleResult(result regionResult) {
	newResult := clientResult{
		Error: result.Error,
		Data:  result.Data,
	}
	middle := middleResult{
		result: newResult,
		times:  tableQueue.getRegionNumber(result.TableName),
	}
	clientQueue[clientQueue.Find(result.ClientIP)].resultQueue <- middle
}

func handleError(result_e regionResult) {
	newResult := clientResult{
		Error: result_e.Error,
		Data:  result_e.Data,
	}
	middle := middleResult{
		result: newResult,
		times:  2,
	}
	clientQueue[clientQueue.Find(result_e.ClientIP)].resultQueue <- middle
}
