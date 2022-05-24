package master

type clientInfo struct {
	ipAddress   string
	resultQueue chan clientResult
}

type clientList []*clientInfo

func (cq clientList) Len() int { return len(cq) }

func (cq clientList) Find(ipAddress string) int {
	index := 0
	for ; index < cq.Len(); index++ {
		if cq[index].ipAddress == ipAddress {
			return index
		}
	}
	return -1
}

func (cq clientList) Remove(ipAddress string) {
	index := cq.Find(ipAddress)
	cq = append(cq[:index], cq[index+1:]...)
}
