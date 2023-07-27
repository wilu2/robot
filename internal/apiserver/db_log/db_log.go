package dblog

import (
	"context"
	"financial_statement/internal/apiserver/dal"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/pkg/log"
	"sync"
	"time"
)

var (
	once sync.Once
	db   dal.Repo
)

// 写db日志
func DbLog(ctx context.Context, taskId uint32, msg string) {
	once.Do(func() {
		db, _ = dal.GetDbFactoryOr(nil)
	})
	var (
		tLog = query.Use(db.GetDb()).Log
	)

	log.Infof("Task:%d DB Log:%s", taskId, msg)
	tLog.WithContext(ctx).Create(&model.Log{
		TaskID:    taskId,
		Msg:       msg,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	})
}
