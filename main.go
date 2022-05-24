package main

import (
	"fmt"
	"sync"
	"log"
	"my/AlienDB/master"
	"time"
)

var lock sync.Mutex

func main() {
	fmt.Println(GetLocalIP())
	go ListenClient()
	go ListenRegion()
	for {}
}