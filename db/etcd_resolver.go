package db

import (
	"context"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"sync"
)

type etcdResolver struct {
	ctx        context.Context
	cancel     context.CancelFunc
	cc         resolver.ClientConn
	etcdClient *clientv3.Client
	scheme     string
	ipPool     sync.Map
}

func (e *etcdResolver) ResolveNow(options resolver.ResolveNowOptions) {
	logrus.Println("etcd resolver resolve now")
}

func (e *etcdResolver) Close() {
	logrus.Println("etcd resolver close")
	e.cancel()
}

func (e *etcdResolver) watcher() {
	watchChan := e.etcdClient.Watch(context.Background(), "/"+e.scheme, clientv3.WithPrefix())

	for {
		select {
		case val := <-watchChan:
			for _, event := range val.Events {
				switch event.Type {
				case 0: // 有数据增加
					e.store(event.Kv.Key, event.Kv.Value)
					logrus.Println("put: ", string(event.Kv.Key))
					e.updateState()
				case 1: // 有数据减少
					logrus.Println("del: ", string(event.Kv.Key))
					e.updateState()
				}
			}
		case <-e.ctx.Done():
			return
		}
	}
}

func (e *etcdResolver) store(k, v []byte) {
	e.ipPool.Store(string(k), string(v))
}

func (e *etcdResolver) del(key []byte) {
	e.ipPool.Delete(string(key))
}

func (e *etcdResolver) updateState() {
	var addrlist resolver.State

	e.ipPool.Range(func(key, value any) bool {
		tA, ok := value.(string)
		if !ok {
			return false
		}
		logrus.Printf("conn.UpdateState key[%v]; val[%v]", key, value)
		addrlist.Addresses = append(addrlist.Addresses, resolver.Address{
			Addr: tA,
		})
		return true
	})

	e.cc.UpdateState(addrlist)
}
