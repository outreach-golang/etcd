package etcd

import (
	"go.etcd.io/etcd/clientv3"
	"time"
)

//NewEtcdClient etcd客户端创建
func NewEtcdClient(endpoints []string) (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return cli, nil
}
