package etcdm

import (
	"context"
	"fmt"
	"time"

	"github.com/404sec/log"

	uuid "github.com/satori/go.uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

type etcdClient struct {
	Client *clientv3.Client
}

func (s *etcdClient) Close() error {
	return s.Client.Close()
}

func (s *etcdClient) GetClient() *clientv3.Client {
	return s.Client
}
func (s *etcdClient) Register(ctx context.Context, prefix, serviceName, addr string) error {
	// 创建一个租约
	client := s.Client
	lease := clientv3.NewLease(client)
	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	leaseResp, err := lease.Grant(cancelCtx, 3)
	if err != nil {
		return err
	}

	leaseChannel, err := lease.KeepAlive(ctx, leaseResp.ID) // 长链接, 不用设置超时时间
	if err != nil {
		return err
	}

	em, err := endpoints.NewManager(client, prefix)
	if err != nil {
		return err
	}

	cancelCtx, cancel = context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	if err := em.AddEndpoint(cancelCtx, fmt.Sprintf("%s/%s/%s", prefix, serviceName, uuid.NewV4().String()), endpoints.Endpoint{
		Addr: addr,
	}, clientv3.WithLease(leaseResp.ID)); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case resp := <-leaseChannel:
				if resp != nil {
					//log.Println("keep alive success.")
				} else {
					log.Infow(ctx, "Etcd Register", "keep alive failed.")
					time.Sleep(time.Second)
					continue
				}
			case <-ctx.Done():
				log.Infow(ctx, "Etcd Register", "close service register")
				cancelCtx, cancel = context.WithTimeout(ctx, time.Second*3)
				defer cancel()
				err := em.DeleteEndpoint(cancelCtx, serviceName)
				if err != nil {
					log.Errorw(ctx, "err", err.Error())
				}
				err = lease.Close()
				if err != nil {
					log.Errorw(ctx, "err", err.Error())
				}
				err = client.Close()
				if err != nil {
					log.Errorw(ctx, "err", err.Error())
				}
				return
			}
		}
	}()

	return nil
}
