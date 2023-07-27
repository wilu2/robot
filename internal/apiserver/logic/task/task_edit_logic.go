package task

import (
	"context"
	"encoding/json"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/consts"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	dblog "financial_statement/internal/apiserver/db_log"
	taskresultasync "financial_statement/internal/apiserver/job/task-result-async"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/task"
	taskhelper "financial_statement/internal/pkg/task_helper"
	"financial_statement/pkg/errors"
	"financial_statement/pkg/log"
	"fmt"

	"github.com/gin-gonic/gin"
)

type TaskEditLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewTaskEditLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) TaskEditLogic {
	return TaskEditLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// TaskEdit 编辑财报页面
func (l *TaskEditLogic) TaskEdit(req *task.EditTaskReq) (err error) {
	var (
		u     = l.ginCtx.Keys["user"].(*model.User)
		tTask = query.Use(l.svcCtx.Db).Task
		q     = query.Use(l.svcCtx.Db)
	)

	financialStatement, err := json.Marshal(req.FinancialStatement)
	if err != nil {
		err = errors.WithCodeMsg(code.Internal, err.Error())
		return
	}

	sqlQuery := tTask.WithContext(l.ctx).Where(tTask.ID.Eq(req.TaskID))
	if l.svcCtx.Config.ServerOptions.IsSaas {
		sqlQuery = sqlQuery.Where(tTask.CreaterUserID.Eq(u.ID))
	}
	dbTask, err := sqlQuery.First()
	if err != nil {
		err = errors.WithCodeMsg(code.Internal, err.Error())
		return
	}

	if err = q.Transaction(func(tx *query.Query) error {
		var path string
		if path, err = taskhelper.SaveTaskResult(dbTask, financialStatement); err != nil {
			return err
		}
		if _, err = tx.Task.WithContext(l.ctx).Where(tTask.ID.Eq(req.TaskID)).Update(tTask.StandardResult, path); err != nil {
			err = errors.WithCodeMsg(code.Internal, err.Error())
			return err
		}
		if req.CancellationFlag {
			if _, err = tx.Task.WithContext(l.ctx).Where(tTask.ID.Eq(req.TaskID)).Update(tTask.Status, consts.TaskCancellationFlag); err != nil {
				err = errors.WithCodeMsg(code.Internal, err.Error())
				return err
			}
		}
		return nil
	}); err != nil {
		err = errors.WithCodeMsg(code.Internal, err.Error())
		return err
	}

	// 发布一个任务结果异步回调处理任务
	if _, err := taskresultasync.NewAsyncTask(&taskresultasync.SyncTaskInfo{
		TaskId:       req.TaskID,
		TaskType:     taskresultasync.SyncTaskTypeSaved,
		CurrentIndex: req.CurrentIndex,
	}); err != nil {
		log.Errorf("new async task with err:%s", err.Error())
	}
	dblog.DbLog(l.ctx, req.TaskID, fmt.Sprintf("save task by user：%s", u.Account))
	return
}
