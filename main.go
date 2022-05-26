package main

import (
	"client"
	"sync"
)

var lock sync.Mutex

func main() {
	//fmt.Println(master_manage.GetLocalIP())
	//go master_manage.ListenClient()
	//go master_manage.ListenRegion()
	client.RunClient()
	for {
	}
}
