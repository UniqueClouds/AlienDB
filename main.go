package main

import (
	"log"
	"my/AlienDB/master"
	"time"
)

func main() {
	//client.RunClient()
	//return
	func() {
		var endPoints = []string{"localhost:2379"}
		ser := master.NewServiceDiscover(endPoints)
		defer ser.Close()
		_ = ser.WatchService()
		//ser.WatchService("/gRPC/")
		for {
			select {
			case <-time.Tick(10 * time.Second):
				log.Println(ser.GetServices())
			}
		}
	}()

	//master.RunMaster()
	//client.RunClient()
	//client.Interpreter("create table fuck;")
	//client.Test()
}
