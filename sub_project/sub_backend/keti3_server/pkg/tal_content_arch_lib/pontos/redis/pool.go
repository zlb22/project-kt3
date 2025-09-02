package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	*redis.Client
}

type RedisConfig struct {
	Addr         string        `mapstructure:"addr"`
	Auth         string        `mapstructure:"auth"`
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

func (r *RedisConfig) initConfig() {
	if r.PoolSize == 0 {
		r.PoolSize = 10
	}
	if r.MinIdleConns == 0 {
		r.MinIdleConns = 5
	}
	if r.IdleTimeout == 0 {
		r.IdleTimeout = 300 * time.Second
	}
}

func NewClient(conf *RedisConfig) (*Client, error) {
	conf.initConfig()
	option := &redis.Options{
		Addr:         conf.Addr,
		Password:     conf.Auth,
		DB:           conf.DB,
		PoolSize:     conf.PoolSize,
		MinIdleConns: conf.MinIdleConns,
		IdleTimeout:  conf.IdleTimeout,
	}

	client := redis.NewClient(option)

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	return &Client{client}, nil
}
