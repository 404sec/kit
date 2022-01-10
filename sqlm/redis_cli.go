package sqlm

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
)

type RedisDBConfig struct {
	Addr         string
	Password     string
	DB           int
	PoolSize     int
	ReadTimeout  int
	DialTimeout  int
	MinIdleConns int //Minimum number of idle connections which is useful when establishing new connection is slow.
}

func (r *RedisDBConfig) newConfig() (conf *redis.Options) {

	conf = &redis.Options{
		Addr:         r.Addr,
		Password:     r.Password,
		DB:           r.DB,
		PoolSize:     r.PoolSize,
		ReadTimeout:  time.Duration(r.ReadTimeout) * time.Second,
		DialTimeout:  time.Duration(r.DialTimeout) * time.Second,
		MinIdleConns: r.MinIdleConns, //Minimum number of idle connections which is useful when establishing new connection is slow.

	}

	return conf
}

//
func (r *RedisDBConfig) NewConnect(ctx context.Context) (client *redis.Client, err error) {

	client = redis.NewClient(r.newConfig())
	client.WithContext(ctx)
	_, err = client.Ping(context.Background()).Result()
	return client, err
}

/**/
