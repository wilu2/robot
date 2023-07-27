package dal

import (
	genericOptions "financial_statement/internal/pkg/options"
	"financial_statement/pkg/storage"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
)

type RedisRepo interface {
	GetRedisDb() redis.UniversalClient
	CloseRedisDb() error

	GetAsynqRedis() *asynq.Client
	GetAsynqRedisConnOpt() *asynq.RedisConnOpt
}

type redisRepo struct {
	Db                *redis.UniversalClient
	AsynqRedis        *asynq.Client
	AsynqRedisConnOpt *asynq.RedisConnOpt
}

func (r *redisRepo) GetAsynqRedisConnOpt() *asynq.RedisConnOpt {
	return r.AsynqRedisConnOpt
}

func (r *redisRepo) GetAsynqRedis() *asynq.Client {
	return r.AsynqRedis
}

func (r *redisRepo) GetRedisDb() redis.UniversalClient {
	return *r.Db
}

func (r *redisRepo) CloseRedisDb() error {
	db := r.GetRedisDb()
	return db.Close()
}

var (
	redisFactory RedisRepo
	redisOnce    sync.Once
)

func GetRedisFactoryOr(opts *genericOptions.RedisOptions) (RedisRepo, error) {
	if opts == nil && redisFactory == nil {
		return nil, fmt.Errorf("failed to get redis store fatory")
	}
	var err error
	var redisClientIns *redis.UniversalClient
	var asynqClientIns *asynq.Client
	var asynqRedisConnOpt *asynq.RedisConnOpt
	redisOnce.Do(func() {
		options := &storage.RedisOptions{
			Type:               opts.Type,
			MasterName:         opts.MasterName,
			SentinelAddrs:      opts.SentinelAddrs,
			SentinelUsername:   opts.SentinelUsername,
			SentinelPassword:   opts.SentinelPassword,
			ClusterAddrs:       opts.ClusterAddrs,
			Host:               opts.Host,
			Port:               opts.Port,
			Username:           opts.Username,
			Password:           opts.Password,
			Database:           opts.Database,
			PoolSize:           opts.PoolSize,
			MinIdleConns:       opts.MinIdleConns,
			DialTimeOut:        opts.DialTimeOut,
			ReadTimeout:        opts.ReadTimeOut,
			WriteTimeout:       opts.WriteTimeOut,
			PoolTimeout:        opts.PoolTimeOut,
			IdleCheckFrequency: opts.IdleCheckFreq,
			IdleTimeout:        opts.IdleTimeout,
		}
		redisClientIns, err = storage.RedisNew(options)
		if err != nil {
			return
		}
		asynqClientIns, asynqRedisConnOpt, err = storage.AsynqRedisNew(options)
		if err != nil {
			return
		}
		redisFactory = &redisRepo{
			Db:                redisClientIns,
			AsynqRedis:        asynqClientIns,
			AsynqRedisConnOpt: asynqRedisConnOpt,
		}
	})
	if err != nil {
		return redisFactory, fmt.Errorf("failed to get redis db with error:%s", err.Error())
	}
	return redisFactory, nil
}
