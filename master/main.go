package master

import (
	"fmt"
)

var regionList PriorityQueue

func RunMaster() {
	fmt.Println("Master 运行中 ...")
	ListenRegion()
}
