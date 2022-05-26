package main

import (
	"fmt"
	master "my/AlienDB/master/master_manage"
)

func main() {
	fmt.Println("The ip address of Master-Server:", master.GetLocalIP())
	go master.RunServiceDiscovery()
	go master.ListenClient()
	go master.ListenRegion()
	for {
	}
}
