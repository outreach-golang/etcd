package etcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

func (s *ServiceRegister) WatchKvChage(key string, channels chan<- string) {
	var (
		getResponse   *clientv3.GetResponse
		watcher       clientv3.Watcher
		watchRespChan <-chan clientv3.WatchResponse
		watchResp     clientv3.WatchResponse
		event         *clientv3.Event
		err           error
	)
	if getResponse, err = s.cli.Get(context.TODO(), key); err != nil {
		fmt.Println(err)
		return
	}
	//创建watcher
	watcher = clientv3.NewWatcher(s.cli)

	// 开始监听， 从当前获取版本的下一个版本开始监听
	watchRespChan = watcher.Watch(context.TODO(), key, clientv3.WithRev(getResponse.Header.Revision+1))

	// 处理kv变化事件
	for watchResp = range watchRespChan {
		for _, event = range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为:", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
				channels <- string(event.Kv.Value)
			case mvccpb.DELETE:
				fmt.Println("删除了", "Revision:", event.Kv.ModRevision)
			}
		}
	}
}

func (s *ServiceRegister) PullValue(key string) {

}
