package master

import (
	"container/heap"
	"fmt"
)

type IpAddressInfo struct {
	ipAddress   string // The value of the item; arbitrary.
	tableNumber int    // The priority of the item in the queue.
	index       int    // The index of the item in the heap.
	ch 			chan string
}

type PriorityQueue []*IpAddressInfo

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].tableNumber > pq[j].tableNumber
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
			break
		}
	}
	return index
}

func (pq *PriorityQueue) getNextTwo() (string, string) {
	old := *pq
	firstLast := old[0]
	secondLast := old[1]
	// old.update(firstLast, firstLast.ipAddress, firstLast.tableNumber-1)
	// old.update(secondLast, secondLast.ipAddress, secondLast.tableNumber-1)
	return firstLast.ipAddress, secondLast.ipAddress
}

func (pq *PriorityQueue) remove(ipAddress string) {
	index := 0
	for ; index < pq.Len(); index++ {
		if ipAddress == (*pq)[index].ipAddress {
			break
		}
	}
	heap.Remove(pq, index)
}

func (pq *PriorityQueue) update(ipAddressInfo *IpAddressInfo, ipAddress string, tableNumber int) {
	ipAddressInfo.ipAddress = ipAddress
	ipAddressInfo.tableNumber = tableNumber
	heap.Fix(pq, ipAddressInfo.index)
}

func test() {
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

	item := &IpAddressInfo{
		ipAddress:   "orange",
		tableNumber: 1,
	}
	heap.Push(&pq, item)
	pq.update(item, item.ipAddress, 5)

	first, second := pq.getNextTwo()
	fmt.Printf("%s, %s ", first, second)

	first, second = pq.getNextTwo()
	fmt.Printf("%s, %s ", first, second)

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*IpAddressInfo)
		fmt.Printf("%s:%d ", item.ipAddress, item.tableNumber)
	}
}
