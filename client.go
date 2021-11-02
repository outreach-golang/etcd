package etcd

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/coreos/etcd/clientv3"
	"io/ioutil"
	"sync"
	"time"
)

var GEtcd *Etcd

type Etcd struct {
	Conf *Config
	init sync.Once
	err  error
	cli  *clientv3.Client
}

func init() {
	GEtcd = newEtcd()
}

func newEtcd() *Etcd {
	return &Etcd{init: sync.Once{}}
}

func (e *Etcd) InitEtcd(ops ...Option) error {
	e.init.Do(func() {

		conf := DefaultConf()

		for _, op := range ops {
			op(conf)
		}

		e.Conf = conf

		var tlsConfig = &tls.Config{}

		if conf.env != "default" && conf.needSSL == 1 {

			cert, errLkp := tls.LoadX509KeyPair(conf.dirPath+conf.sslServerFile, conf.dirPath+conf.sslKeyFile)
			if e.err = errLkp; e.err != nil {
				return
			}

			caData, errRf := ioutil.ReadFile(conf.dirPath + conf.caFile)
			if e.err = errRf; e.err != nil {
				return
			}

			pool := x509.NewCertPool()
			pool.AppendCertsFromPEM(caData)

			tlsConfig.Certificates = []tls.Certificate{cert}
			tlsConfig.RootCAs = pool
		} else {
			tlsConfig = nil
		}

		e.cli, e.err = newEtcdClient(conf.points, tlsConfig)
	})

	return e.err
}

func (e *Etcd) GetCli() *clientv3.Client {
	return e.cli
}

//newEtcdClient etcd客户端创建
func newEtcdClient(endpoints []string, config *tls.Config) (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
		TLS:         config,
	})
	if err != nil {
		return nil, err
	}

	return cli, nil
}
