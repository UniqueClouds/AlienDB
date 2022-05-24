package main

var tableQueue tableList
var clientQueue clientList
var regionQueue PriorityQueue

func handleCreate(ipAddress, tableName, sqlInstruction string) {
	desRegion_1, desRegion_2 := regionQueue.getNextTwo()
	request := regionRequest {
		TableName: tableName,
		IpAddress: ipAddress,
		Kind: "create",
		Sql: sqlInstruction,
		File: nil,
	}
	regionQueue[regionQueue.find(desRegion_1)].requestQueue <- request
	regionQueue[regionQueue.find(desRegion_2)].requestQueue <- request
}

func handleOther(ipAddress, tableName, kind, sqlInstruction string) {
	desRegion_1, desRegion_2 := tableQueue.getRegionIp(tableName)
	request := regionRequest {
		TableName: tableName,
		IpAddress: ipAddress,
		Kind: kind,
		Sql: sqlInstruction,
		File: nil,
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

func handleResult(result regionResult){
	newResult := clientResult {
		Error: result.Error,
		Data: result.Data,
	}
	clientQueue[clientQueue.Find(result.ClientIP)].resultQueue <- newResult
}

func handleError(result_e regionResult){
	newResult := clientResult {
		Error: result_e.Error,
		Data: result_e.Data,
	}
	clientQueue[clientQueue.Find(result_e.ClientIP)].resultQueue <- newResult
}