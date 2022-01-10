package sqlm

import (
	"github.com/go-redis/redis/v8"
)

type RedisDB struct {
	Client *redis.Client
	Cfg    *RedisDBConfig
}

type RedisDBConfig struct {
	Addr         string
	Password     string
	DB           int
	PoolSize     int
	ReadTimeout  int
	DialTimeout  int
	MinIdleConns int //Minimum number of idle connections which is useful when establishing new connection is slow.
}
