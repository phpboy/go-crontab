package core

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var etcd *etcdObj

type etcdObj struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
}

func InitEtcd() {

	var (
		client *clientv3.Client
		err error
		config clientv3.Config
		kv clientv3.KV
		lease clientv3.Lease
	)

	config=clientv3.Config{
		Endpoints:[]string{Config.Etcd},
		DialTimeout:5*time.Second,
	}

	if client,err = clientv3.New(config);err!=nil{
		fmt.Println(err)
	}

	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	etcd = &etcdObj{
		client:client,
		kv:kv,
		lease:lease,
	}
}
