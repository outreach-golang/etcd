package etcd

import (
	"context"
	"errors"
	"os"
	"sync"
)

var (
	ServiceRegisterHandler *ServiceInfo
)

type ServiceInfo struct {
	serviceName string
	nodeIP      string
	path        string
	lease       int64
	init        sync.Once
	err         error
}

func init() {
	ServiceRegisterHandler = newService()
}

func newService() *ServiceInfo {
	return &ServiceInfo{
		nodeIP: "http://" + os.Getenv("NODE_IP"),
		path:   "exterior-gateway/services-conf/services/",
		lease:  30,
		init:   sync.Once{},
		err:    nil,
	}
}

func (s *ServiceInfo) InitServiceRegister(ctx context.Context, sr *ServiceRegister, serviceName string,
	port string,
) error {
	s.init.Do(func() {
		if _, has := os.LookupEnv("NODE_IP"); !has {

			s.err = errors.New("环境变量 NODE_IP 获取失败！")

			return
		}

		if serviceName == "" {
			s.err = errors.New("注册服务时 ServiceName 不能为空！")

			return
		}

		serviceName = s.path + serviceName + "." + getRandomString(12)

		accessAddress := s.nodeIP + ":" + port

		if err := sr.Register(ctx, serviceName, accessAddress, s.lease); err != nil {
			s.err = err
		}
	})

	return s.err
}
