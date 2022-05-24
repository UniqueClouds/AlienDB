package main

import (
	"fmt"
	"sync"
)

var lock sync.Mutex

func main() {
	fmt.Println("The ip address of Master-Server:", GetLocalIP())
	go RunServiceDiscovery()
	go ListenClient()
	go ListenRegion()
	for {}
}