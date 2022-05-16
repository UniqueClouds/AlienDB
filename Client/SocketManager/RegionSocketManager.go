package SocketManager

// reference: Go Socket Programming
// https://www.cnblogs.com/Yunya-Cnblogs/p/13815864.html
import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"net"
	"os"
	"sync"
	"tonydb/RegionServer/minisql/src/CatalogManager"
)
var socket
var host = flag.String("localhost", "", "localhost")
var port = flag.String("port", "3306", "port")
var region string = "localhost"
func setRegionIP(ip string) {
	region = ip
}

func connectRegionServer(PORT int)  {
	flag.Parse()
	conn, err := net.Dial("tcp", *host+":"+*port)
	if err != nil{
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connecting to "+ *host + ":" + *port)

	done := make(chan string)

	var wg sync.WaitGroup
	wg.Add(2)7

	go handleWrite(conn, &wg)
	go handleRead(conn, &wg)

	fmt.Println(<-done)
	fmt.Println(<-done)

	wg.Wait()
}
func handleRead(conn net.Conn, wg *sync.WaitGroup){
	defer wg.Done()
	buf := make([]byte, 1024)
	reqLen, err := conn.Read(buf)
	if err != nil{
		fmt.Println("Error to read message because of", err)
		return
	}

}
func handleWrite(conn net.Conn,wg *sync.WaitGroup ){
	defer wg.Done()

}

