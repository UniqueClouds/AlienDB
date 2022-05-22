package master

import (
	"container/heap"
	"fmt"
)

var ipQueue PriorityQueue

func fakeSample() {
	ipAddressInfos := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}
	ipQueue = make(PriorityQueue, len(ipAddressInfos))
	i := 0
	for ipAddress, tableNumber := range ipAddressInfos {
		ipQueue[i] = &IpAddressInfo{
			ipAddress:   ipAddress,
			tableNumber: tableNumber,
			index:       i,
			ch:     	 make(chan string, 20),
		}
		i++
	}
	heap.Init(&ipQueue)
	fmt.Println(ipQueue.Len())
	handleCreate("hh", "haha")
}

func handleCreate(tableName, sqlInstruction string) {
	desRegion_1, desRegion_2 := ipQueue.getNextTwo()
	ipQueue[ipQueue.find(desRegion_1)].ch <- sqlInstruction
	ipQueue[ipQueue.find(desRegion_2)].ch <- sqlInstruction
}

func handleOther(tableName, sqlInstruction string) {
	// desRegion_1, desRegion_2 := ipQueue.getNextTwo()
	
}

func forwardToRegion(sqlInstruction string) {
	
}