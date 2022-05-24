package main

import (
	"fmt"
	master "my/AlienDB/master/master_manage"

	"sync"
)

var lock sync.Mutex

func main() {
	fmt.Println(master.GetLocalIP())
	go master.RunServiceDiscovery()
	go master.ListenClient()
	go master.ListenRegion()
	for {
	}
}
