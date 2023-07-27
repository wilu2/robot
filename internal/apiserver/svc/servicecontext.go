package svc

import (
	"financial_statement/internal/apiserver/config"
	"financial_statement/internal/apiserver/dal"
	"financial_statement/internal/apiserver/job"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config  *config.Config
	Db      *gorm.DB
	RedisDb redis.UniversalClient
	JobMgr  *job.JobMgr
}

func NewServiceContext(cfg *config.Config, jobMgr *job.JobMgr) *ServiceContext {
	dbIns, _ := dal.GetDbFactoryOr(nil)
	redisIns, _ := dal.GetRedisFactoryOr(nil)
	return &ServiceContext{
		Config:  cfg,
		Db:      dbIns.GetDb(),
		RedisDb: redisIns.GetRedisDb(),
		JobMgr:  jobMgr,
	}
}
