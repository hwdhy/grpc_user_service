package db

import (
	"context"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type EtcdRegister struct {
	etcdCli *clientv3.Client // etcd连接
	leaseId clientv3.LeaseID // 租约ID
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewEtcdRegister() (*EtcdRegister, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:12379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logrus.Printf("new etcd client faild, error %v", err)
		return nil, err
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	svr := &EtcdRegister{
		etcdCli: client,
		ctx:     ctx,
		cancel:  cancelFunc,
	}
	return svr, nil
}

// CreateLease 创建租约
func (e *EtcdRegister) CreateLease(expire int64) error {
	res, err := e.etcdCli.Grant(e.ctx, expire)
	if err != nil {
		return err
	}
	e.leaseId = res.ID
	return nil
}

// BindLease 绑定租约
func (e *EtcdRegister) BindLease(key string, value string) error {
	res, err := e.etcdCli.Put(e.ctx, key, value, clientv3.WithLease(e.leaseId))
	if err != nil {
		return err
	}

	logrus.Printf("bindLease success %v ", res)
	return nil
}

// KeepAlive 租约 发送心跳，表名服务正常
func (e *EtcdRegister) KeepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	resChan, err := e.etcdCli.KeepAlive(e.ctx, e.leaseId)
	if err != nil {
		logrus.Printf("keepAlive failed,error %v ", err)
	}
	return resChan, err
}

func (e *EtcdRegister) Watcher(key string, resChan <-chan *clientv3.LeaseKeepAliveResponse) {
	for {
		select {
		case l := <-resChan:
			logrus.Printf("续约成功, val: %+v ", l)
		case <-e.ctx.Done():
			logrus.Printf("续约关闭")
			return
		}
	}
}

func (e *EtcdRegister) Close() error {
	e.cancel()

	logrus.Printf("close...")
	// 撤销租约
	e.etcdCli.Revoke(e.ctx, e.leaseId)
	return e.etcdCli.Close()
}

// RegisterServer 注册服务
func (e *EtcdRegister) RegisterServer(serviceName string, addr string, expire int64) error {
	err := e.CreateLease(expire)
	if err != nil {
		return err
	}

	err = e.BindLease(serviceName, addr)
	if err != nil {
		return nil
	}

	keepAliveChan, err := e.KeepAlive()
	if err != nil {
		return err
	}

	go e.Watcher(serviceName, keepAliveChan)
	return nil
}
