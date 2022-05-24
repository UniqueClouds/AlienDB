package main

import (
	"fmt"
	"sync"
)

var lock sync.Mutex

func main() {
	fmt.Println(GetLocalIP())
	go ListenClient()
	go ListenRegion()
	for {}
}