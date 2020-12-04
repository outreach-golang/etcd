package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

func (s *ServiceRegister) WatchKvChange(key string, channels chan<- string) {
	var (
		getResponse   *clientv3.GetResponse
		watcher       clientv3.Watcher
		watchRespChan <-chan clientv3.WatchResponse
		watchResp     clientv3.WatchResponse
		event         *clientv3.Event
		err           error
	)

	if getResponse, err = s.Get(key); err != nil {
		fmt.Println("WatchKvChange get the key fail", err)
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

func (s *ServiceRegister) Get(key string) (*clientv3.GetResponse, error) {
	var (
		getResp *clientv3.GetResponse
		err     error
	)
	if getResp, err = s.kv.Get(context.TODO(), key); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return getResp, nil
}

func (s *ServiceRegister) PullValue(key, name string) (interface{}, bool) {
	var (
		getResponse *clientv3.GetResponse
		data        []byte
		values      map[string]interface{}
		err         error
	)
	if getResponse, err = s.Get(key); err != nil {
		fmt.Println("WatchKvChange get the key fail", err)
	}

	//没有做匹配查询，所有默认只会取出一个值
	data = getResponse.Kvs[0].Value

	values = s.JsonToMap(data)

	if v, ok := values[name]; ok {
		return v, ok
	} else {
		return "", false
	}

}

func (s *ServiceRegister) JsonToMap(data []byte) map[string]interface{} {
	m := make(map[string]interface{})

	fmt.Println(string(data))

	err := json.Unmarshal(data, &m)
	if err != nil {
		fmt.Printf("Unmarshal with error: %+v\n", err)
		return map[string]interface{}{}
	}
	return m
}
