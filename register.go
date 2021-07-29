package etcd

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	jsoniter "github.com/json-iterator/go"
)

//ServiceRegister 创建租约注册服务
type ServiceRegister struct {
	cli     *clientv3.Client
	leaseID clientv3.LeaseID

	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
	val           string
	ctx           context.Context
}

//NewServiceRegister 新建注册服务
func NewServiceRegister(cli *clientv3.Client) *ServiceRegister {
	sr := &ServiceRegister{
		cli: cli,
	}

	return sr
}

//Register 注册服务
func (s *ServiceRegister) Register(cxt context.Context, key string, val string, lease int64) error {

	s.key = key
	s.val = val
	s.ctx = cxt

	//和已经注册的服务合并
	sameService, _ := s.getTheSameService(cxt, key)
	sameService.Ips = append(sameService.Ips, val)
	//去掉重复的服务
	sameService.Ips = removeDuplicationMap(sameService.Ips)

	s.val, _ = jsoniter.MarshalToString(&sameService)

	if err := s.putKeyWithLease(lease); err != nil {

		return err
	}

	return nil
}

type registerService struct {
	Ips []string `json:"ips"`
}

//检查是否已经存在相同服务
func (s *ServiceRegister) getTheSameService(cxt context.Context, key string) (registerService, error) {
	rs := registerService{}

	response, err := s.cli.KV.Get(cxt, key)
	if err != nil {
		return rs, err
	}

	for _, kv := range response.Kvs {
		_ = jsoniter.Unmarshal(kv.Value, &rs)
	}

	return rs, nil
}

//设置租约
func (s *ServiceRegister) putKeyWithLease(lease int64) error {

	resp, err := s.cli.Grant(context.Background(), lease)
	if err != nil {
		return err
	}

	_, err = s.cli.Put(context.Background(), s.key, s.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	leaseRespChan, err := s.cli.KeepAlive(context.Background(), resp.ID)

	if err != nil {
		return err
	}
	s.leaseID = resp.ID

	s.keepAliveChan = leaseRespChan

	return nil
}

//ListenLeaseRespChan 监听 续租情况
func (s *ServiceRegister) ListenLeaseRespChan() *clientv3.LeaseKeepAliveResponse {
	for leaseKeepResp := range s.keepAliveChan {
		return leaseKeepResp

	}
	return nil
}

// Close 注销服务
func (s *ServiceRegister) Close() error {

	if _, err := s.cli.Revoke(context.Background(), s.leaseID); err != nil {
		return err
	}

	return s.cli.Close()
}

//DeleteKey 删除Key
func (s *ServiceRegister) DeleteKey(key string) error {
	_, err := s.cli.Delete(context.Background(), key)
	if err != nil {
		return err
	}

	return nil
}
