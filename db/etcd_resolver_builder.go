package db

import (
	"context"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
)

type EtcdResolverBuilder struct {
	etcdClient *clientv3.Client
}

func NewEtcdResolverBuilder() *EtcdResolverBuilder {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:12379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal("client get etcd failed, error: ", err)
	}
	return &EtcdResolverBuilder{
		etcdClient: etcdClient,
	}
}

func (e *EtcdResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	prefix := "/" + target.URL.Scheme
	logrus.Println(prefix)

	res, err := e.etcdClient.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	es := &etcdResolver{
		cc:         cc,
		etcdClient: e.etcdClient,
		ctx:        ctx,
		cancel:     cancelFunc,
		scheme:     target.URL.Scheme,
	}
	logrus.Printf("etcd res: %+v", res)
	for _, kv := range res.Kvs {
		es.store(kv.Key, kv.Value)
	}

	es.updateState()

	go es.watcher()
	return es, nil
}

func (e *EtcdResolverBuilder) Scheme() string {
	return "etcd"
}
