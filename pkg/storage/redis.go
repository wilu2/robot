package storage

import (
	"context"
	"financial_statement/pkg/log"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
)

var (
	RedisTypeDefault  = "default"
	RedisTypeSentinel = "sentinel"
	RedisTypeCluster  = "cluster"
)

// RedisOptions Config defines options for redis cluster.
type RedisOptions struct {
	Type               string
	MasterName         string
	SentinelAddrs      string
	SentinelUsername   string
	SentinelPassword   string
	ClusterAddrs       string
	Host               string
	Port               int
	Username           string
	Password           string
	Database           int
	PoolSize           int
	MinIdleConns       int
	DialTimeOut        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	PoolTimeout        time.Duration
	IdleCheckFrequency time.Duration
	IdleTimeout        time.Duration
}

func RedisNew(opts *RedisOptions) (*redis.UniversalClient, error) {
	var redisClient redis.UniversalClient
	log.Infof("Redis 初始化：%v , 模式：%s", opts, opts.Type)
	var universalClientOpts = &redis.UniversalOptions{
		DB:                 opts.Database,
		Username:           opts.Username,
		Password:           opts.Password,
		PoolSize:           opts.PoolSize,
		MinIdleConns:       opts.MinIdleConns,
		DialTimeout:        opts.DialTimeOut,
		ReadTimeout:        opts.ReadTimeout,
		WriteTimeout:       opts.WriteTimeout,
		PoolTimeout:        opts.PoolTimeout,
		IdleCheckFrequency: opts.IdleCheckFrequency,
		IdleTimeout:        opts.IdleTimeout,
	}
	switch opts.Type {
	case RedisTypeDefault:
		// 单节点客户端
		universalClientOpts.Addrs = []string{fmt.Sprintf("%s:%d", opts.Host, opts.Port)}
	case RedisTypeSentinel:
		// 哨兵客户端
		if strings.Contains(opts.SentinelAddrs, ";") {
			universalClientOpts.Addrs = strings.Split(opts.SentinelAddrs, ";")
		} else {
			universalClientOpts.Addrs = strings.Split(opts.SentinelAddrs, ",")
		}
		universalClientOpts.MasterName = opts.MasterName
		universalClientOpts.SentinelUsername = opts.SentinelUsername
		universalClientOpts.SentinelPassword = opts.SentinelPassword
	case RedisTypeCluster:
		// 集群客户端
		if strings.Contains(opts.ClusterAddrs, ";") {
			universalClientOpts.Addrs = strings.Split(opts.ClusterAddrs, ";")
		} else {
			universalClientOpts.Addrs = strings.Split(opts.ClusterAddrs, ",")
		}
		universalClientOpts.RouteRandomly = true
	}
	redisClient = redis.NewUniversalClient(universalClientOpts)

	timeoutCtx, cancel := context.WithTimeout(context.Background(), opts.PoolTimeout*time.Second)
	defer cancel()
	if _, err := redisClient.Ping(timeoutCtx).Result(); err != nil {
		return nil, fmt.Errorf("failed to get redis db")
	}
	return &redisClient, nil
}

func AsynqRedisNew(opts *RedisOptions) (*asynq.Client, *asynq.RedisConnOpt, error) {
	var opt asynq.RedisConnOpt
	log.Infof("Asynq Redis 初始化：%v , 模式：%s", opts, opts.Type)
	if opts.Type == RedisTypeDefault {
		// 单节点客户端
		opt = &asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%d", opts.Host, opts.Port), Password: opts.Password, DB: opts.Database}
	} else if opts.Type == RedisTypeSentinel {
		var addrs []string
		if strings.Contains(opts.SentinelAddrs, ";") {
			addrs = strings.Split(opts.SentinelAddrs, ";")
		} else {
			addrs = strings.Split(opts.SentinelAddrs, ",")
		}
		// 哨兵客户端
		opt = &asynq.RedisFailoverClientOpt{
			MasterName:       opts.MasterName,
			SentinelAddrs:    addrs,
			SentinelPassword: opts.SentinelPassword,
			Username:         opts.Username,
			Password:         opts.Password,
			DB:               opts.Database,
			DialTimeout:      opts.DialTimeOut,
			ReadTimeout:      opts.ReadTimeout,
			WriteTimeout:     opts.WriteTimeout,
			PoolSize:         opts.PoolSize,
		}
	} else if opts.Type == RedisTypeCluster {
		var addrs []string
		if strings.Contains(opts.ClusterAddrs, ";") {
			addrs = strings.Split(opts.ClusterAddrs, ";")
		} else {
			addrs = strings.Split(opts.ClusterAddrs, ",")
		}
		// 集群客户端
		opt = &asynq.RedisClusterClientOpt{
			Addrs:        addrs,
			Username:     opts.Username,
			Password:     opts.Password,
			DialTimeout:  opts.DialTimeOut,
			ReadTimeout:  opts.ReadTimeout,
			WriteTimeout: opts.WriteTimeout,
		}
	}
	client := asynq.NewClient(opt)
	return client, &opt, nil
}
