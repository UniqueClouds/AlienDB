package master

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"sync"
	"time"
)

var copyInfoChannel chan []copyTableInfo

type ServiceDiscovery struct {
	cli        *clientv3.Client
	serverList map[string]string // 服务列表
	lock       sync.RWMutex
}

// NewServiceDiscover 新建发现服务
func NewServiceDiscover(endpoints []string) *ServiceDiscovery {
	// 初始化
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		log.Fatal(err)
	}

	return &ServiceDiscovery{
		cli:        cli,
		serverList: make(map[string]string),
	}
}

//WatchService 初始化服务列表和监视
func (s *ServiceDiscovery) WatchService() error {
	prefix := "/db/"
	// 根据前缀获取现有的 key
	resp, err := s.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	// 遍历获取到的 key 和 value
	for _, ev := range resp.Kvs {
		s.SetServiceList(string(ev.Key), string(ev.Value))
	}

	// 监视前缀， 修改变更的 server
	go s.watcher()
	return nil
}

//watcher 监听key的前缀
func (s *ServiceDiscovery) watcher() {
	prefix := "/db/"
	rch := s.cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
	log.Printf("watching prefix:%s now ...", prefix)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT: // 修改或者新增
				s.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE: // 删除
				ip := string(ev.Kv.Value)
				s.DelServiceList(string(ev.Kv.Key))
				copyInfoChannel <- tableQueue.downRegionIp(ip)
			}
		}
	}
}

//SetServiceList 新增服务地址
func (s *ServiceDiscovery) SetServiceList(key, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.serverList[key] = string(val)
	fmt.Println(">>> etcd: put key :", key, "val", val)
	log.Println("put key :", key, "val:", val)
}

//DelServiceList 删除服务地址
func (s *ServiceDiscovery) DelServiceList(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.serverList, key)
	fmt.Println(">>> etcd: del key", key)
	log.Println("del key", key)
}

//GetServices 获取服务地址
func (s *ServiceDiscovery) GetServices() []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	addrs := make([]string, 0)

	for _, v := range s.serverList {
		addrs = append(addrs, v)
	}
	return addrs
}

//Close 关闭服务
func (s *ServiceDiscovery) Close() error {
	return s.cli.Close()
}

func RunServiceDiscovery() {
	var endPoints = []string{"localhost:2379"}
	ser := NewServiceDiscover(endPoints)
	defer ser.Close()
	ser.WatchService()
	for {
		select {
		case <-time.Tick(10 * time.Second):
			log.Println(ser.GetServices())
		}
	}
}

//func test() {
//	var endPoints = []string{"localhost:2379"}
//	ser := NewServiceDiscover(endPoints)
//	defer ser.Close()
//	ser.WatchService()
//	//ser.WatchService("/gRPC/")
//	for {
//		select {
//		case <-time.Tick(10 * time.Second):
//			log.Println(ser.GetServices())
//		}
//	}
//}
