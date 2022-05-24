package main

import (
	"container/heap"
	"fmt"
	"sort"
)

type IpAddressInfo struct {
	ipAddress   	string // The value of the item; arbitrary.
	tableNumber 	int    // The priority of the item in the queue.
	requestQueue 	chan regionRequest
	receiveQueue    chan regionResult
	index       	int    // The index of the item in the heap.
}

type PriorityQueue []*IpAddressInfo

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].tableNumber < pq[j].tableNumber
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	ipAddressInfo := x.(*IpAddressInfo)
	ipAddressInfo.index = n
	*pq = append(*pq, ipAddressInfo)
	sort.Slice(*pq, pq.Less)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	ipAddressInfo := old[n-1]
	old[n-1] = nil           // avoid memory leak
	ipAddressInfo.index = -1 // for safety
	*pq = old[0 : n-1]
	return ipAddressInfo
}

func (pq *PriorityQueue) find(ipAddress string) (int) {
	index := 0
	for ; index < pq.Len(); index++ {
		if ipAddress == (*pq)[index].ipAddress {
			return index
		}
	}
	return -1
}

func (pq *PriorityQueue) getNextTwo() (string, string) {
	old := *pq
	firstLast := old[0]
	secondLast := old[1]
	// old.update(firstLast, firstLast.ipAddress, firstLast.tableNumber-1)
	// old.update(secondLast, secondLast.ipAddress, secondLast.tableNumber-1)
	return firstLast.ipAddress, secondLast.ipAddress
}

func (pq *PriorityQueue) getCopyRegion(aliveIp string) (string) {
	if (*pq)[0].ipAddress != aliveIp {
		return (*pq)[0].ipAddress;
	}
	return (*pq)[1].ipAddress
}

func removeRegion(pq PriorityQueue, ipAddress string) PriorityQueue {
	index := 0
	for ; index < pq.Len(); index++ {
		if ipAddress == pq[index].ipAddress {
			break
		}
	}
	return append(pq[:index], pq[index+1:]...)
}

func (pq *PriorityQueue) update(ipAddress string, tableNumber int) {
	ipAddressInfo := (*pq)[pq.find(ipAddress)]
	ipAddressInfo.tableNumber = tableNumber
	sort.Slice(*pq, pq.Less)
}

func test_1() {
	ipAddressInfos := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}

	pq := make(PriorityQueue, len(ipAddressInfos))
	i := 0
	for ipAddress, tableNumber := range ipAddressInfos {
		pq[i] = &IpAddressInfo{
			ipAddress:   ipAddress,
			tableNumber: tableNumber,
			index:       i,
		}
		i++
	}
	heap.Init(&pq)

	sort.Slice(pq, pq.Less)
	newI := &IpAddressInfo {
		ipAddress: "hh",
		tableNumber: 0,
	}

	pq.Push(newI)
	pq.update("hh", 9)
	pq = removeRegion(pq, "apple")
	pq = removeRegion(pq, "banana")
	for i := 0; i < pq.Len(); i++ {
		fmt.Printf("%s:%d ", pq[i].ipAddress, pq[i].tableNumber)
	}

	fmt.Println(pq.getCopyRegion("hh"))
	fmt.Println(pq.getCopyRegion("pear"))
}
