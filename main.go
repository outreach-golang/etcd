package etcd

import (
	"context"
	"github.com/outreach-golang/logger"
	"go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
	"time"
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
func NewServiceRegister(cxt context.Context, endpoints []string, key, val string, lease int64) (*ServiceRegister, error) {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logger.WithContext(cxt).Fatal(err.Error())
	}

	ser := &ServiceRegister{
		cli: cli,
		key: "server/" + key,
		val: val,
		ctx: cxt,
	}

	if err := ser.putKeyWithLease(lease); err != nil {

		logger.WithContext(ser.ctx).Fatal(err.Error())

		return nil, err
	}

	return ser, nil
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

	logger.WithContext(s.ctx).Info("leaseID", zap.Int64("leaseID", int64(resp.ID)))

	s.keepAliveChan = leaseRespChan

	logger.WithContext(s.ctx).Info("Put key success!", zap.String("key", s.key), zap.String("val", s.val))

	return nil
}

//ListenLeaseRespChan 监听 续租情况
func (s *ServiceRegister) ListenLeaseRespChan() {
	for leaseKeepResp := range s.keepAliveChan {
		logger.WithContext(s.ctx).Info("续约成功", zap.Any("leaseKeepResp", leaseKeepResp))
	}
	logger.WithContext(s.ctx).Info("关闭续租")
}

// Close 注销服务
func (s *ServiceRegister) Close() error {

	if _, err := s.cli.Revoke(context.Background(), s.leaseID); err != nil {
		return err
	}

	logger.WithContext(s.ctx).Info("撤销租约")

	return s.cli.Close()
}
