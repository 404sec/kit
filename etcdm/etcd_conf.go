package etcdm

import (
	"context"
	"crypto/tls"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type EtcdConf struct {
	EndPoints []string
	TLSInfo   TLSInfo
}

// TLSInfo etcd中证书的相关配置
type TLSInfo struct {
	CertFile      string
	KeyFile       string
	TrustedCAFile string
}

func (s *EtcdConf) NewConfig() (conf clientv3.Config, err error) {
	conf = clientv3.Config{
		Endpoints:   s.EndPoints,
		DialTimeout: 5 * time.Second,
	}
	if s.TLSInfo.CertFile != "" &&
		s.TLSInfo.KeyFile != "" &&
		s.TLSInfo.TrustedCAFile != "" {
		tlsInfo := transport.TLSInfo{
			CertFile:      s.TLSInfo.CertFile,
			KeyFile:       s.TLSInfo.KeyFile,
			TrustedCAFile: s.TLSInfo.TrustedCAFile,
		}
		var tlsConfig *tls.Config
		tlsConfig, err = tlsInfo.ClientConfig()
		if err != nil {
			return conf, err
		}
		conf.TLS = tlsConfig
	}

	return conf, nil
}

func (s *EtcdConf) NewConnect(ctx context.Context) (etcdClient, error) {

	var cl etcdClient
	econf, err := s.NewConfig()
	if err != nil {
		return cl, err
	}
	cl.Client, err = clientv3.New(econf)
	return cl, err
}
