package etcd

import (
	"context"
	"go.etcd.io/etcd/clientv3"
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
func (s *ServiceRegister) NewServiceRegister(cli *clientv3.Client) error {
	s.cli = cli

	return nil
}

//Register 注册服务
func (s *ServiceRegister) Register(cxt context.Context, key string, val string, lease int64) error {

	s.key = key
	s.val = val
	s.ctx = cxt

	if err := s.putKeyWithLease(lease); err != nil {

		return err
	}

	return nil
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
