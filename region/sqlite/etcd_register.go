package sqlite

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// reference: https://bingjian-zhu.github.io/2020/05/14/etcd%E5%AE%9E%E7%8E%B0%E6%9C%8D%E5%8A%A1%E5%8F%91%E7%8E%B0/

// ServiceRegister 服务注册
type ServiceRegister struct {
	cli     *clientv3.Client // etcd client
	leaseID clientv3.LeaseID // 租约ID
	// 租约keepalive 相应 chan
	keepAlveChan <-chan *clientv3.LeaseKeepAliveResponse
	key          string // key
	val          string // value
}

// NewServiceRegister 创建租约注册服务
func NewServiceRegister(endpoints []string, key, val string, lease int64, dialTimeout int) (*ServiceRegister, error) {
	fmt.Println(">Region: 向主节点注册服务中 ....")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Duration(dialTimeout) * time.Second,
	})

	if err != nil {
		return nil, err
	}

	ser := &ServiceRegister{
		cli: cli,
		key: key,
		val: val,
	}

	if err := ser.putKeyWithLease(lease); err != nil {
		return nil, err
	}
	return ser, nil
}

// 设置租约
func (s *ServiceRegister) putKeyWithLease(lease int64) error {
	// 创建一个新的租约
	resp, err := s.cli.Grant(context.Background(), lease)
	if err != nil {
		return err
	}

	// 注册服务并绑定租约
	_, err = s.cli.Put(context.Background(), s.key, s.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}
	// 设置租约 定期发送需求请求
	// keepalive 是给定的租约永远有效 如果发布到通道的keepalive响应没有被立即使用
	// 则租约客户端至少每秒钟向etcd服务器发送保持活动请求
	// etcd client 会自动发送ttl到etcd server，从而保证该租约一直有效
	leaseRespChan, err := s.cli.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return err
	}
	s.leaseID = resp.ID
	log.Println(s.leaseID)
	s.keepAlveChan = leaseRespChan
	log.Printf("Put key: %s val: %s success", s.key, s.val)
	return nil
}

// ListenLeaseRespChan 监听 续租情况
func (s *ServiceRegister) ListenLeaseRespChan() {
	for leaseKeepResp := range s.keepAlveChan {
		log.Println(">Region: 续约成功", leaseKeepResp)
	}
	log.Println("关闭续租")
}

// Close 注销服务
func (s *ServiceRegister) Close() error {
	if _, err := s.cli.Revoke(context.Background(), s.leaseID); err != nil {
		return err
	}
	log.Println("撤销租约")
	return s.cli.Close()
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// RegionRegister 从节点服务注册功能
func RegionRegister(localAddr string) {
	var endpoints = []string{"192.168.81.166:2379"}
	//var endpoints = []string{"localhost:2379"}
	//localIP := GetLocalIP()
	// 暂定名称

	ser, err := NewServiceRegister(endpoints, "/db/region"+RandStringBytes(4), localAddr, 6, 5)
	if err != nil {
		log.Fatalln(err)
	}
	// 监听续租相应 chan
	go ser.ListenLeaseRespChan()

	for {

	}
}
