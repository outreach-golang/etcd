package etcd

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"sync"
)

type WatchPutFn func(key, val string)

type WatchDelFn func(key, val string)

type Kv struct {
	cli    *clientv3.Client
	KvList map[string]string
	lock   sync.RWMutex
}

func NewKv(cli *clientv3.Client) *Kv {
	return &Kv{
		cli:    cli,
		KvList: make(map[string]string),
		lock:   sync.RWMutex{},
	}
}

//WatchKeyByPrefix 根据前缀监听kv
func (k *Kv) WatchKeyByPrefix(ctx context.Context, prefix string, putFn, delFn WatchPutFn) error {

	response, err := k.cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, kv := range response.Kvs {
		k.setKv(string(kv.Key), string(kv.Value))
	}

	go k.watchKeyByPrefix(ctx, prefix, putFn, delFn)

	return nil

}

func (k *Kv) watchKeyByPrefix(ctx context.Context, prefix string, putFn, delFn WatchPutFn) {
	watchChan := k.cli.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())

	for response := range watchChan {
		for _, event := range response.Events {
			switch event.Type {
			case mvccpb.PUT:
				putFn(string(event.Kv.Key), string(event.Kv.Value))
				break
			case mvccpb.DELETE:
				delFn(string(event.Kv.Key), string(event.PrevKv.Value))
				break
			}
		}
	}

}

//setKv 设置kv
func (k *Kv) setKv(key, val string) {
	k.lock.Lock()
	defer k.lock.Unlock()

	k.KvList[key] = val
}

//delKv 删除kv
func (k *Kv) delKv(key string) {
	k.lock.Lock()
	defer k.lock.Unlock()

	delete(k.KvList, key)
}

//Close 关闭客户端
func (k *Kv) Close() error {
	return k.cli.Close()
}
