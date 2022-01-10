package sqlm

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
)

func (r *RedisDB) newConfig() (conf *redis.Options, err error) {

	conf = &redis.Options{
		Addr:         r.Cfg.Addr,
		Password:     r.Cfg.Password,
		DB:           r.Cfg.DB,
		PoolSize:     r.Cfg.PoolSize,
		ReadTimeout:  time.Duration(r.Cfg.ReadTimeout) * time.Second,
		DialTimeout:  time.Duration(r.Cfg.DialTimeout) * time.Second,
		MinIdleConns: r.Cfg.MinIdleConns, //Minimum number of idle connections which is useful when establishing new connection is slow.

	}

	return conf, nil
}

//
func (r *RedisDB) NewConnect(ctx context.Context) (err error) {
	conf, err := r.newConfig()
	if err != nil {
		return err
	}
	r.Client = redis.NewClient(conf)
	r.Client.WithContext(ctx)
	_, err = r.Client.Ping(context.Background()).Result()
	return err
}

func (r *RedisDB) Close() error {
	return r.Client.Close()
}

//
func (r *RedisDB) GetClient() *redis.Client {
	return r.Client
}

/**/
